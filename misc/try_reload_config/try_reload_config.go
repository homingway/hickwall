package main

import (
	"fmt"
	"github.com/davecheney/profile"
	"github.com/oliveagle/hickwall/backends"
	"github.com/oliveagle/hickwall/collectors"
	"github.com/oliveagle/hickwall/config"
	"github.com/oliveagle/hickwall/lib"
	"github.com/oliveagle/hickwall/utils"
	"time"
)

func main() {
	_ = collectors.GetCollectors

	config.LoadRuntimeConfFromFileOnce()

	fmt.Println("--")
	fmt.Println(backends.GetBackendList())

	runtime_conf := config.GetRuntimeConf()
	fmt.Println(runtime_conf.Transport_stdout)
	fmt.Println(runtime_conf.Transport_file)
	fmt.Println(runtime_conf.Transport_influxdb)

	runtime_conf.Transport_file.Path = "./filedump.txt"
	runtime_conf.Transport_file.Enabled = true
	runtime_conf.Transport_stdout.Enabled = false

	runtime_conf.Transport_influxdb[0].Enabled = false
	runtime_conf.Transport_influxdb[1].Enabled = false
	config.UpdateRuntimeConf(runtime_conf)

	backends.CreateBackendsFromRuntimeConf()
	backends.RunBackends()

	tick := time.Tick(time.Millisecond * 5000)
	done := time.After(time.Minute * 8)

	collectors.CreateCollectorsFromRuntimeConf()
	collectors.RunCollectors()

	utils.HttpPprofServe(6060)

	cfg := profile.Config{
		CPUProfile:     true,
		MemProfile:     true,
		BlockProfile:   true,
		ProfilePath:    "./pprofs/", // store profiles in current directory
		NoShutdownHook: true,        // do not hook SIGINT
	}
	defer profile.Start(&cfg).Stop()

loop:
	for {
		select {
		case <-tick:
			lib.ReloadWithRuntimeConfig()
		case md := <-collectors.GetDataChan():
			// for _, p := range md {
			// 	fmt.Println(" point : ", p)
			// }
			backends.WriteToBackends(md)
		case <-done:
			break loop
		}
	}

}
