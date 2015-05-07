package main

import (
	"fmt"
	"github.com/oliveagle/hickwall/backends"
	"github.com/oliveagle/hickwall/collectorlib"
	// . "github.com/oliveagle/hickwall/collectors"
	"github.com/oliveagle/hickwall/config"
	"math/rand"
	"time"
)

func main() {

	config.LoadRuntimeConfFromFileOnce()
	backends.CreateBackendsFromRuntimeConf()

	fmt.Println("--")
	fmt.Println(backends.GetBackendList())

	runtime_conf := config.GetRuntimeConf()
	fmt.Println(runtime_conf.Transport_stdout)
	fmt.Println(runtime_conf.Transport_file)
	runtime_conf.Transport_stdout.Enabled = true
	config.UpdateRuntimeConf(runtime_conf)

	file_bk, _ := backends.GetBackendByName("file")
	go file_bk.Run()
	defer file_bk.Close()

	tick := time.Tick(time.Millisecond * 100)
	done := time.After(time.Second * 6)

	// TODO: file backend donesn't work
loop:
	for {
		select {
		case <-tick:
			// backend.Write(md)
			rand.Seed(time.Now().UnixNano())
			p := collectorlib.DataPoint{
				Metric:    "metric1",
				Timestamp: time.Now(),
				Value:     rand.Intn(100),
			}
			md := collectorlib.MultiDataPoint{p}

			file_bk.Write(md)
		case <-done:
			break loop
		}
	}

}
