package mem_test

import (
	// "log"
	"github.com/oliveagle/hickwall/config"
	"github.com/oliveagle/hickwall/hickwall"
	"path/filepath"
	"testing"
	"time"
)

func Test_Mem(t *testing.T) {
	config.CONF_FILEPATH, _ = filepath.Abs("../../hickwall/test/config_mem.yml")
	t.Log(config.CONF_FILEPATH)
	core, _, err := hickwall.LoadConfigStrategyFile()
	if err != nil {
		t.Error("failed")
		return
	}
	defer core.Close()

	// tick := time.Tick(time.Second * 1)
	// tickClose := time.Tick(time.Second * 2)
	done := time.After(time.Second * 3)

	// first_value := 0
	// last_value := 0
	// delta := 0

	for {
		select {
		// case <-tickClose:
		// 	pc.Close()
		// 	pc = cgo_pdh.NewPdhCollector()
		// 	pc.AddCounter("\\Process(hickwall)\\Working Set - Private")
		// 	log.Println("close and recreate pdh collector")
		// case <-tick:
		// 	data, err = pc.CollectAllDouble()
		// 	for _, d := range data {
		// 		if first_value == 0 {
		// 			first_value = int(d) / 1024
		// 		}
		// 		last_value = int(d) / 1024
		// 		delta = last_value - first_value
		// 		log.Printf("first: %d, last: %d, delta: %d\n", first_value, last_value, delta)
		// 	}
		case <-done:
			return
		}
	}
}
