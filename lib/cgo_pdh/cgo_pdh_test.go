package cgo_pdh

import (
	"fmt"
	"math"
	"testing"
	"time"
)

var (
	_ = fmt.Sprintf("")
	_ = time.Now()
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

the following test replaced to .misc/lib_cgo_pdh_test to avoid jamming test_all.sh
*/
// func Test_cgo_pdh_mem_leak_multiple_instance(t *testing.T) {
// ...
// }
