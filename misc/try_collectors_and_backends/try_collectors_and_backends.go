package main

import (
	"fmt"
	"github.com/kr/pretty"
	"github.com/oliveagle/hickwall/backends"
	"github.com/oliveagle/hickwall/collectorlib"
	"github.com/oliveagle/hickwall/collectors"
	// "github.com/oliveagle/hickwall/config"
	"time"

	"github.com/oliveagle/hickwall/utils"
)

func main() {
	utils.HttpPprofServe(6060)

	pretty.Println("")

	fmt.Println(backends.GetBackendList())

	ch := make(chan collectorlib.MultiDataPoint)

	collectors.RunCollectors(ch)

	backends.RunBackends()
	defer backends.CloseBackends()

	done := time.After(time.Second * 600)
	delay := time.After(time.Second * 1)
loop:
	for {
		select {
		case md, err := <-ch:
			fmt.Println("MultiDataPoint: ", md, err)
			for _, p := range md {
				fmt.Println(" point ---> ", p)
			}
			backends.WriteToBackends(md)
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
