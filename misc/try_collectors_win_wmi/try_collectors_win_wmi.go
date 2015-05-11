package main

import (
	"fmt"
	"github.com/oliveagle/hickwall/collectorlib"
	"github.com/oliveagle/hickwall/collectors"
	"github.com/oliveagle/hickwall/config"
	"time"

	"github.com/kr/pretty"
)

func main() {
	pretty.Println("")
	config.LoadRuntimeConfFromFileOnce()

	runtime_conf := config.GetRuntimeConf()
	fmt.Println(runtime_conf.Collector_win_wmi)

	// collectors.AddCollector("win_wmi", "cc[0]collector", runtime_conf.Collector_win_wmi[0])
	collectors.AddCollector("win_wmi", "cc[0]collector", runtime_conf.Collector_win_wmi)
	cc := collectors.GetCollectors()

	fmt.Println(" ++ collectors:  ", cc)

	ch := make(chan collectorlib.MultiDataPoint)

	collectors.RunCollectors(ch)

	done := time.After(time.Second * 6)
loop:
	for {
		select {
		case dp, err := <-ch:
			fmt.Println("MultiDataPoint: ", dp, err)
			// case <-ch:
			// fmt.Println(" point ---> ", dp, err)
			// fmt.Println("-------------------")
			// pretty.Println(dp)
			for _, p := range dp {
				fmt.Println(" point ---> ", p)
			}
		case <-done:
			break loop
		}
	}
}
