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

	// test set hostname from runtime config
	runtime_conf := config.GetRuntimeConf()
	runtime_conf.Hostname = "tes哈哈啊哈t"
	config.UpdateRuntimeConf(runtime_conf)

	cs := GetBuiltinCollectors()
	fmt.Println("runtime_conf.Collector_win_pdh: ", runtime_conf.Collector_win_pdh)

	AddCustomizedCollectorByName("win_pdh", "cc[0]collector", runtime_conf.Collector_win_pdh)
	// AddCustomizedCollectorByName("win_pdh", "cc[1]collector", runtime_conf.Collector_win_pdh[1])
	cc := GetCustomizedCollectors()

	fmt.Println(" ++ customized_collectors:  ", cc)
	fmt.Println(" ++ builtin_collectors: ", cs)

	for _, c := range cs {
		fmt.Println(c.Name())
	}

	ch := make(chan collectorlib.MultiDataPoint)

	RunBuiltinCollectors(ch)

	go cs[0].Run(ch)

	// go cc[0].Run(ch)
	// go cc[1].Run(ch)

	done := time.After(time.Second * 30)
	delay := time.After(time.Second * 5)
loop:
	for {
		select {
		case dp, err := <-ch:
			fmt.Println("MultiDataPoint: ----------------------------------", err)
			// case <-ch:
			// fmt.Println(" point ---> ", dp, err)
			// fmt.Println("-------------------")
			// pretty.Println(dp)
			for _, p := range dp {
				fmt.Println(" point ---> ", p)
			}
		case <-delay:

			// cs[0].Stop()
			// StopBuiltinCollectors()

			// StopCustomizedCollectors()
			// fmt.Println("customized_collectors", GetCustomizedCollectors())
			// RemoveAllCustomizedCollectors()
			// fmt.Println("customized_collectors", GetCustomizedCollectors())

			// AddCustomizedCollectorByName("win_pdh", "cc[0]collector", runtime_conf.Collector_win_pdh[0])
			// AddCustomizedCollectorByName("win_pdh", "cc[1]collector", runtime_conf.Collector_win_pdh[1])
			// fmt.Println("customized_collectors", GetCustomizedCollectors())

			// RunCustomizedCollectors(ch)

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
