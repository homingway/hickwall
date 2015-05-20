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

	AddCollector("cmd", "cc[0]collector", runtime_conf.Collector_cmd)
	cc := GetCollectors()

	fmt.Println(" ++ collectors:  ", cc)

	ch := make(chan collectorlib.MultiDataPoint)

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
