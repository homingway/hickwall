package main

import (
	"fmt"
	"github.com/oliveagle/hickwall/collectorlib"
	. "github.com/oliveagle/hickwall/collectors"
	"github.com/oliveagle/hickwall/config"
	"time"

	"github.com/kr/pretty"
)

func main() {
	pretty.Println("")
	config.LoadRuntimeConfFromFileOnce()
	runtime_conf := config.GetRuntimeConf()

	fmt.Println("runtime_conf.Collector_win_pdh: ", runtime_conf.Collector_win_sys)
	AddCollector("sys_windows", "cc[0]collector", runtime_conf.Collector_win_sys)
	cc := GetCollectors()

	fmt.Println(" ++ collectors:  ", cc)

	for _, c := range cc {
		fmt.Println(c.Name())
	}

	ch := make(chan collectorlib.MultiDataPoint)

	RunCollectors(ch)

	// go cc[0].Run(ch)
	// go cc[1].Run(ch)

	done := time.After(time.Second * 30)
	delay := time.After(time.Second * 5)
loop:
	for {
		select {
		case dp, err := <-ch:
			fmt.Println("MultiDataPoint: ----------------------------------", err)
			for _, p := range dp {
				fmt.Println(" point ---> ", p)
			}
		case <-delay:

			// StopCollectors()
			// fmt.Println("collectors", GetCollectors())
			// RemoveAllCollectors()
			// fmt.Println("collectors", GetCollectors())

			// AddCollector("win_pdh", "cc[0]collector", runtime_conf.Collector_win_pdh[0])
			// AddCollector("win_pdh", "cc[1]collector", runtime_conf.Collector_win_pdh[1])
			// fmt.Println("collectors", GetCollectors())

			// RunCollectors(ch)

			// cc[0].Stop()
			// cc[1].Stop()

			// change config on the fly
			// cs[0].Init()
			// cs[0].(*IntervalCollector).SetInterval(time.Millisecond * 200)
		case <-done:
			break loop
		}
	}
}
