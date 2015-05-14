package main

import (
	"fmt"
	// "github.com/davecheney/profile"
	"github.com/oliveagle/hickwall/backends"
	"github.com/oliveagle/hickwall/collectorlib"
	"github.com/oliveagle/hickwall/collectors"
	"github.com/oliveagle/hickwall/config"
	"math/rand"
	// "os"
	"runtime"
	"runtime/debug"
	// "runtime/pprof"
	"time"
)

func main() {
	_ = collectors.GetCollectors
	_ = runtime.GOROOT()
	// runtime.MemProfileRate
	debug.SetGCPercent(50)

	config.LoadRuntimeConfFromFileOnce()

	fmt.Println("--")
	fmt.Println(backends.GetBackendList())

	runtime_conf := config.GetRuntimeConf()
	fmt.Println(runtime_conf.Transport_stdout)
	fmt.Println(runtime_conf.Transport_file)

	runtime_conf.Transport_file.Path = "./filedump.txt"
	runtime_conf.Transport_file.Enabled = true
	runtime_conf.Transport_stdout.Enabled = true

	fmt.Println(runtime_conf.Transport_stdout)
	fmt.Println(runtime_conf.Transport_file)
	config.UpdateRuntimeConf(runtime_conf)

	backends.CreateBackendsFromRuntimeConf()

	file_bk, _ := backends.GetBackendByName("file")
	go file_bk.Run()
	defer file_bk.Close()

	tick := time.Tick(time.Millisecond * 1000)
	// done := time.After(time.Second * 6)
	done := time.After(time.Second * 130)

	collectors.AddCollector("hickwall_client", "hickwall_client", nil)
	for _, c := range collectors.GetCollectors() {
		fmt.Println(c)
	}
	collectors.RunCollectors()

	// cfg := profile.Config{
	// 	CPUProfile:     true,
	// 	MemProfile:     true,
	// 	BlockProfile:   true,
	// 	ProfilePath:    "./pprofs/", // store profiles in current directory
	// 	NoShutdownHook: true,        // do not hook SIGINT
	// }
	// defer profile.Start(&cfg).Stop()

	// utils.HttpPprofServe(6060)

	// first := time.After(time.Second * 10)
	// second := time.After(time.Minute * 2)

	// TODO: file backend donesn't work
loop:
	for {
		select {
		case md := <-collectors.GetDataChan():
			// fmt.Println(" point: ", md[0])
			// for _, p := range md {
			// 	fmt.Println(" point : ", p)
			// }

			file_bk.Write(md)

		case <-tick:
			// no memory leak in this branch !!! file backend is fine.
			rand.Seed(time.Now().UnixNano())
			p := &collectorlib.DataPoint{
				Metric:    "metric1",
				Timestamp: time.Now(),
				Value:     rand.Intn(100),
			}
			md := collectorlib.MultiDataPoint{p}
			file_bk.Write(md)
			debug.FreeOSMemory()

		// case <-first:
		// 	f, err := os.Create("./mem.pprof.1")
		// 	if err != nil {
		// 		fmt.Println(err)
		// 		break loop
		// 	}
		// 	runtime.GC()
		// 	pprof.WriteHeapProfile(f)
		// 	f.Close()
		// case <-second:
		// 	f, err := os.Create("./mem.pprof.2")
		// 	if err != nil {
		// 		fmt.Println(err)
		// 		break loop
		// 	}
		// 	runtime.GC()
		// 	pprof.WriteHeapProfile(f)
		// 	f.Close()
		case <-done:
			break loop
		}
	}

}
