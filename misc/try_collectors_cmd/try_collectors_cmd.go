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
	// backends.CreateBackendsFromRuntimeConf()
	runtime_conf := config.GetRuntimeConf()

	cs := GetBuiltinCollectors()

	AddCustomizedCollectorByName("cmd", "cc[0]collector", runtime_conf.Collector_cmd[0])
	cc := GetCustomizedCollectors()

	fmt.Println(" ++ customized_collectors:  ", cc)
	fmt.Println(" ++ builtin_collectors: ", cs)

	ch := make(chan collectorlib.MultiDataPoint)

	go cs[0].Run(ch)
	go cc[0].Run(ch)

	done := time.After(time.Second * 3)
	delay := time.After(time.Second * 1)
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
		case <-delay:
			// change config on the fly
			// cs[0].Init()
			// cs[0].(*IntervalCollector).SetInterval(time.Millisecond * 200)
		case <-done:
			break loop
		}
	}
}
