package pdh

import (
	"fmt"
	"testing"
	"time"
)

var (
	_ = time.Now()
	_ = fmt.Sprintln("")
)

func TestPdh(t *testing.T) {
	pc := NewPdhCollector()
	pc.AddEnglishCounter("\\System\\Processes")
	data := pc.CollectData()
	defer pc.Close()

	for _, d := range data {
		t.Log(d.Value)
		if d.Err != nil || d.Value < 10 {
			t.Error("processes count less than 10 ?", d.Err)
		}
	}
}

func TestPdhAddInvalidCounterPath(t *testing.T) {
	pc := NewPdhCollector()
	pc.AddEnglishCounter("\\Systemhjahahah\\Processes")
	data := pc.CollectData()
	defer pc.Close()

	for _, d := range data {
		t.Log(d.Value)
		if d.Err == nil || d.Value > 0 {
			t.Error("should emit error")
		}
	}
}

// // confirmed pdh_windows.go don't have memory leak on single instance
// func TestPdh_mem_leak_single_instance(t *testing.T) {
// 	pc := NewPdhCollector()
// 	pc.AddEnglishCounter("\\Process(pdh.test)\\Working Set - Private")
// 	data := pc.CollectData()
// 	defer pc.Close()

// 	tick := time.Tick(time.Second)
// 	done := time.After(time.Second * 100)

// 	first_value := 0
// 	last_value := 0
// 	delta := 0

// 	for {
// 		select {
// 		case <-tick:
// 			for _, d := range data {
// 				if first_value == 0 {
// 					first_value = int(d.Value) / 1024
// 				}
// 				last_value = int(d.Value) / 1024
// 				delta = last_value - first_value

// 				fmt.Printf("first: %d, last: %d, delta: %d\n", first_value, last_value, delta)
// 				if first_value != last_value {
// 					t.Errorf("delta happens:  delta: %d ", delta)
// 					return
// 				}
// 			}
// 		case <-done:
// 			return
// 		}
// 	}
// }

func TestPdh_mem_leak_multiple_instance(t *testing.T) {
	pc := NewPdhCollector()
	pc.AddEnglishCounter("\\Process(pdh.test)\\Working Set - Private")
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
			pc = NewPdhCollector()
			pc.AddEnglishCounter("\\Process(pdh.test)\\Working Set - Private")
		case <-tick:
			for _, d := range pc.CollectData() {
				if first_value == 0 {
					first_value = int(d.Value) / 1024
				}
				last_value = int(d.Value) / 1024
				delta = last_value - first_value

				fmt.Printf("first: %d, last: %d, delta: %d\n", first_value, last_value, delta)
				// if first_value != last_value {
				// 	t.Errorf("delta happens:  delta: %d ", delta)
				// 	return
				// }
			}
		case <-done:
			return
		}
	}
}
