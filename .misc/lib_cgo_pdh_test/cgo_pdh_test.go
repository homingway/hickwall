package lib_cgo_pdh_test

import (
	"fmt"
	"github.com/oliveagle/hickwall/lib/cgo_pdh"
	"testing"
	"time"
)

var (
	_ = fmt.Sprintf("")
	_ = time.Now()
)

/*
place here to avoid jamming test_all
*/
func Test_cgo_pdh_mem_leak_multiple_instance(t *testing.T) {
	pc := cgo_pdh.NewPdhCollector()
	defer pc.Close()
	pc.AddCounter("\\Process(lib_cgo_pdh_test.test)\\Working Set - Private")
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
			pc = cgo_pdh.NewPdhCollector()
			pc.AddCounter("\\Process(lib_cgo_pdh_test.test)\\Working Set - Private")
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
