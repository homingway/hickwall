package main

import (
	"fmt"
	"github.com/kr/pretty"
	"github.com/oliveagle/go-collectors/datapoint"
	. "github.com/oliveagle/hickwall/collectors"
	"time"
)

func main() {
	pretty.Println("")

	cs := GetBuiltinCollectors()
	cc := GetCustomizedCollectors()

	fmt.Println(" ++ builtin_collector: ", cs)
	fmt.Println(" ++ customized_collectors:  ", cc)

	ch := make(chan datapoint.MultiDataPoint)

	for _, c := range cs {
		go c.Run(ch)
	}
	for _, c := range cc {
		go c.Run(ch)
	}

	// go cs.Run(ch)
	// go cc[0].Run(ch)
	// go cc[1].Run(ch)

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
