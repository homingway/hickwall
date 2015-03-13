package main

import (
	"fmt"
	"github.com/kr/pretty"
	"github.com/oliveagle/go-collectors/datapoint"
	"github.com/oliveagle/hickwall/backends"
	"github.com/oliveagle/hickwall/collectors"
	"github.com/oliveagle/hickwall/config"
	"time"

	"github.com/oliveagle/hickwall/utils"
)

func main() {
	utils.HttpPprofServe(6060)

	pretty.Println(config.APP_NAME)

	fmt.Println(backends.GetBackendList())

	ch := make(chan *datapoint.MultiDataPoint)

	collectors.RunAllCollectors(ch)
	backends.RunBackends()
	defer backends.CloseBackends()

	done := time.After(time.Second * 600)
	delay := time.After(time.Second * 1)
loop:
	for {
		select {
		case dp, err := <-ch:
			fmt.Println("MultiDataPoint: ", dp, err)
			for _, p := range *dp {
				fmt.Println(" point ---> ", p)
			}
			backends.WriteToBackends(*dp)
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
