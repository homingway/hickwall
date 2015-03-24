package main

import (
	"fmt"
	"github.com/oliveagle/go-collectors/datapoint"
	. "github.com/oliveagle/hickwall/collectors"
	"time"
)

func main() {

	cs := GetBuiltinCollectorByName("builtin_hickwall_client")
	if cs != nil {

		fmt.Println(" ++ builtin_collectors: ", cs)

		ch := make(chan *datapoint.MultiDataPoint)

		go cs.Run(ch)

		done := time.After(time.Second * 3)
	loop:
		for {
			select {
			case md, err := <-ch:
				fmt.Println("MultiDataPoint: ", md, err)
				for _, p := range *md {
					fmt.Println(" point ---> ", p)
				}
			case <-done:
				break loop
			}
		}
	}
}