package main

import (
	"fmt"
	"github.com/kr/pretty"
	"github.com/oliveagle/hickwall/collectorlib"
	. "github.com/oliveagle/hickwall/collectors"
	"time"
)

func main() {
	pretty.Println("")

	cc := GetCollectors()

	fmt.Println(" ++ collectors:  ", cc)

	ch := make(chan collectorlib.MultiDataPoint)

	for _, c := range cc {
		go c.Run(ch)
	}

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
			// fmt.Println("-------------------")
			// change config on the fly
			// cs[0].Init()
			// cs[0].(*IntervalCollector).SetInterval(time.Millisecond * 200)
		case <-done:
			break loop
		}
	}
}
