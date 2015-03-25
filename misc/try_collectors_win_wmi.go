package main

import (
	"fmt"
	"github.com/oliveagle/go-collectors/datapoint"
	. "github.com/oliveagle/hickwall/collectors"
	"github.com/oliveagle/hickwall/config"
	"time"

	"github.com/kr/pretty"
)

func main() {
	pretty.Println("")

	// pretty.Println(config.Conf)

	cs := GetBuiltinCollectorByName("builtin_win_wmi")

	AddCustomizedCollectorByName("win_wmi", "cc[0]collector", config.Conf.Collector_win_wmi[0])
	cc := GetCustomizedCollectors()

	fmt.Println(" ++ builtin_collector: ", &cs)
	fmt.Println(" ++ customized_collectors:  ", cc)

	ch := make(chan *datapoint.MultiDataPoint)

	// go cs.Run(ch)
	go cc[0].Run(ch)
	// go cc[1].Run(ch)

	done := time.After(time.Second * 60)
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
			for _, p := range *dp {
				fmt.Println(" point ---> ", p)
			}
		case <-delay:
			// fmt.Println("-------------------")
			// change config on the fly
			// cs[0].Init()
			// cs[0].(*IntervalCollector).SetInterval(time.Millisecond * 200)
		case <-done:
			break loop
		}
	}
}
