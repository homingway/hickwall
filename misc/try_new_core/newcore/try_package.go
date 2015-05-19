package main

import (
	"fmt"
	"github.com/oliveagle/hickwall/misc/try_new_core/newcore/backends"
	"github.com/oliveagle/hickwall/misc/try_new_core/newcore/collectors"
	"github.com/oliveagle/hickwall/misc/try_new_core/newcore/newcore"
	"time"
)

var (
	_ = fmt.Sprintf("")
	_ = time.Now()
)

func main() {
	sub := newcore.Subscribe(collectors.NewDummyCollector("c1", time.Millisecond*100), nil)

	// fset := FanOut(sub,
	//  newDummyBackend("b1", time.Second*10),
	//  newDummyBackend("b2", 0))

	fset := newcore.FanOut(sub,
		backends.NewDummyBackend("b1", 0, false))

	fset_closed_chan := make(chan error)

	time.AfterFunc(time.Second*time.Duration(100), func() {
		// sub will be closed within FanOut
		fset_closed_chan <- fset.Close()
	})

	a := 0
	tick := time.Tick(time.Second * time.Duration(1))
	timeout := time.After(time.Second * time.Duration(106))

main_loop:
	for {
		select {
		case <-fset_closed_chan:
			fmt.Println("TestFanout.fset closed")
			break main_loop
		case md, openning := <-sub.Updates():
			if openning == false {
				fmt.Println("TestFanout.sub.Updates() closed")
				break main_loop
			} else {
				fmt.Printf(".")
				// fmt.Printf("TestFanout.sub.Updates() still openning: 0x%X\n", &md)
			}
			a += len(*md)
			// t.Log("md: ", md)
		case <-tick:
			a = 0
		case <-timeout:
			fmt.Println("TestFanout.timed out! something is blocking")
			break main_loop
		}
	}
}
