package cgo_pdh

import (
	"fmt"
	"math"
	"testing"
	"time"
)

// func Test_1(t *testing.T) {
// 	getcpuload()
// }

func Test_2(t *testing.T) {
	p := NewPdhCollector()

	err := p.AddCounter("\\System\\Processes")

	if err != nil {
		t.Error("error: ", err)
	}
	p.AddCounter("\\Process(explorer)\\Working Set - Private")

	res, err := p.CollectAllDouble()
	if err != nil {
		t.Error("error: ", err)
	}
	for idx, value := range res {
		t.Logf("query: %s, value: %f\n", p.GetQueryByIdx(idx), value)
	}
}

func Test_ValidateCStatus(t *testing.T) {
	err := ValidateCStatus(0x0)
	if err != nil {
		t.Error("...")
	}
}

func Test_cgo_pdh_close_invalid_handle(t *testing.T) {
	pc := NewPdhCollector()

	// close multiple times won't have any problem.
	pc.Close()
	pc.Close()
	err := pc.AddCounter("\\Process(cgo_pdh.test)\\Working Set - Private")
	t.Log(err)
	if err == nil {
		t.Error("should fail.")
		return
	}

	data, err := pc.CollectAllDouble()
	t.Log(data, err)
	if err == nil {
		t.Error("error: ", err)
		return
	}
}

func Test_cgo_pdh_collect_invalid_counter(t *testing.T) {
	pc := NewPdhCollector()
	defer pc.Close()
	err := pc.AddCounter("\\Process(cgo_pdh.test)\\Working Set - Private")
	err = pc.AddCounter("\\Process(xxxxxxxxxx)\\Working Set - Private")

	data, err := pc.CollectAllDouble()
	if len(data) != 2 {
		t.Error("should have 2 datapoint")
	}
	if !math.IsNaN(data[1]) {
		t.Error("the second data point should be NaN")
	}
	if err != nil {
		t.Error("invalid counter should not raise error.")
	}

	t.Log(data, err)
}

/* no leaking !!!! great
first: 5888, last: 5872, delta: -16
first: 5888, last: 5872, delta: -16
close and recreate pdh collector
first: 5888, last: 5872, delta: -16
first: 5888, last: 5872, delta: -16
first: 5888, last: 5872, delta: -16
close and recreate pdh collector
first: 5888, last: 5872, delta: -16
close and recreate pdh collector
first: 5888, last: 5872, delta: -16
first: 5888, last: 5872, delta: -16
close and recreate pdh collector
first: 5888, last: 5872, delta: -16
first: 5888, last: 5872, delta: -16
close and recreate pdh collector
first: 5888, last: 5872, delta: -16
first: 5888, last: 5872, delta: -16
close and recreate pdh collector
first: 5888, last: 5872, delta: -16
first: 5888, last: 5872, delta: -16
*/
func Test_cgo_pdh_mem_leak_multiple_instance(t *testing.T) {
	pc := NewPdhCollector()
	defer pc.Close()
	pc.AddCounter("\\Process(cgo_pdh.test)\\Working Set - Private")
	data, err := pc.CollectAllDouble()
	if err != nil {
		t.Error("error: ", err)
		return
	}

	tick := time.Tick(time.Second)
	tickClose := time.Tick(time.Second * 2)
	done := time.After(time.Second * 100)

	first_value := 0
	last_value := 0
	delta := 0

	for {
		select {
		case <-tickClose:
			fmt.Println("close and recreate pdh collector")
			pc.Close()
			pc = NewPdhCollector()
			pc.AddCounter("\\Process(cgo_pdh.test)\\Working Set - Private")
		case <-tick:
			data, err = pc.CollectAllDouble()
			for _, d := range data {
				if first_value == 0 {
					first_value = int(d) / 1024
				}
				last_value = int(d) / 1024
				delta = last_value - first_value

				fmt.Printf("first: %d, last: %d, delta: %d\n", first_value, last_value, delta)
				if delta > 20 {
					t.Error("delta above 20k happened")
				}
			}
		case <-done:
			return
		}
	}
}
