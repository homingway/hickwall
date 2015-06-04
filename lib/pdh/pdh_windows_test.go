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

// the following test moved to .misc/lib_pdh_test to avoid jamming test_all.sh
// // confirmed pdh_windows.go don't have memory leak on single instance
// func TestPdh_mem_leak_single_instance(t *testing.T) {
// ...
// }

// the following test moved to .misc/lib_pdh_test to avoid jamming test_all.sh
// func TestPdh_mem_leak_multiple_instance(t *testing.T) {
// ...
// }
