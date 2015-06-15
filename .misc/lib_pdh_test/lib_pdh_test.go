package lib_pdh_test

import (
	"fmt"
	"github.com/oliveagle/hickwall/lib/pdh"
	"testing"
	"time"
)

// confirmed pdh_windows.go don't have memory leak on single instance
func TestPdh_mem_leak_single_instance(t *testing.T) {
	pc := pdh.NewPdhCollector()
	pc.AddEnglishCounter("\\Process(lib_pdh_test.test)\\Working Set - Private")
	data := pc.CollectData()
	defer pc.Close()

	tick := time.Tick(time.Second)
	done := time.After(time.Second * 100)

	first_value := 0
	last_value := 0
	delta := 0

	for {
		select {
		case <-tick:
			for _, d := range data {
				if first_value == 0 {
					first_value = int(d.Value) / 1024
				}
				last_value = int(d.Value) / 1024
				delta = last_value - first_value

				fmt.Printf("first: %d, last: %d, delta: %d\n", first_value, last_value, delta)
				if first_value != last_value {
					t.Errorf("delta happens:  delta: %d ", delta)
					return
				}
			}
		case <-done:
			return
		}
	}
}

func TestPdh_mem_leak_multiple_instance(t *testing.T) {
	pc := pdh.NewPdhCollector()
	pc.AddEnglishCounter("\\Process(lib_pdh_test.test)\\Working Set - Private")
	// data := pc.CollectData()
	defer pc.Close()

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
			// pc = nil
			pc = pdh.NewPdhCollector()
			pc.AddEnglishCounter("\\Process(lib_pdh_test.test)\\Working Set - Private")
		case <-tick:
			for _, d := range pc.CollectData() {
				if first_value == 0 {
					first_value = int(d.Value) / 1024
				}
				last_value = int(d.Value) / 1024
				delta = last_value - first_value

				fmt.Printf("first: %d, last: %d, delta: %d\n", first_value, last_value, delta)
				// if first_value != last_value {
				//  t.Errorf("delta happens:  delta: %d ", delta)
				//  return
				// }
			}
		case <-done:
			return
		}
	}
}
