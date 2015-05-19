package main

import (
	"fmt"
	"github.com/oliveagle/hickwall/misc/try_new_core/newcore/newcore"
	"runtime"
	"runtime/debug"
	"time"
)

/*
#include "getPeakRSS.c"
*/
import "C"

func main() {
	// Subscribe to some feeds, and create a merged update stream.
	merged := newcore.Merge(
		newcore.Subscribe(newcore.CollectorFactory("c1"), nil),
		newcore.Subscribe(newcore.CollectorFactory("c2"), nil),
		newcore.Subscribe(newcore.CollectorFactory("c3"), nil),
		newcore.Subscribe(newcore.CollectorFactory("c3"), nil),
		newcore.Subscribe(newcore.CollectorFactory("c3"), nil),
		newcore.Subscribe(newcore.CollectorFactory("c1"), nil),
		newcore.Subscribe(newcore.CollectorFactory("c2"), nil),
		newcore.Subscribe(newcore.CollectorFactory("c3"), nil),
		newcore.Subscribe(newcore.CollectorFactory("c3"), nil),
		newcore.Subscribe(newcore.CollectorFactory("c3"), nil),
		newcore.Subscribe(newcore.CollectorFactory("c1"), nil),
		newcore.Subscribe(newcore.CollectorFactory("c2"), nil),
		newcore.Subscribe(newcore.CollectorFactory("c3"), nil),
		newcore.Subscribe(newcore.CollectorFactory("c3"), nil),
		newcore.Subscribe(newcore.CollectorFactory("c3"), nil),
		newcore.Subscribe(newcore.CollectorFactory("c1"), nil),
		newcore.Subscribe(newcore.CollectorFactory("c2"), nil),
		newcore.Subscribe(newcore.CollectorFactory("c3"), nil),
		newcore.Subscribe(newcore.CollectorFactory("c3"), nil),
		newcore.Subscribe(newcore.CollectorFactory("c3"), nil),
		newcore.Subscribe(newcore.CollectorFactory("c1"), nil),
		newcore.Subscribe(newcore.CollectorFactory("c2"), nil),
		newcore.Subscribe(newcore.CollectorFactory("c3"), nil),
		newcore.Subscribe(newcore.CollectorFactory("c3"), nil),
		newcore.Subscribe(newcore.CollectorFactory("c3"), nil),
		newcore.Subscribe(newcore.CollectorFactory("c1"), nil),
		newcore.Subscribe(newcore.CollectorFactory("c2"), nil),
		newcore.Subscribe(newcore.CollectorFactory("c3"), nil),
		newcore.Subscribe(newcore.CollectorFactory("c3"), nil),
		newcore.Subscribe(newcore.CollectorFactory("c3"), nil),
	)

	// Close the subscriptions after some time.
	time.AfterFunc(300*time.Second, func() {
		// fmt.Println("closed:", merged.Close())
		merged.Close()
	})

	var a = 0

	tick := time.Tick(time.Duration(1) * time.Second)
	var mem runtime.MemStats

	fmt.Println("PeakRSS(k), CurrentRSS(k), Alloc(k), Sys(k), HeapSys(k), HeapAlloc(k), HeapInuse(k), HeapIdle(k), HeapReleased(k), HeapObjects, Points/s")

	// var dp *newcore.DataPoint
	var md *newcore.MultiDataPoint
	var channel_closed bool

	debug.SetGCPercent(50)

	for {
		select {
		case md, channel_closed = <-merged.Updates():
			if md == nil && channel_closed == false {
				fmt.Println("merged closed")
				return
			}
			// fmt.Println(md)
			// time.Sleep(time.Millisecond * time.Duration(10))
			a += len(*md)
		case <-tick:
			// debug.FreeOSMemory()
			runtime.ReadMemStats(&mem)
			fmt.Printf("%d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d \n",
				C.getPeakRSS()/1024, C.getCurrentRSS()/1024, mem.Alloc/1024, mem.Sys/1024, mem.HeapSys/1024, mem.HeapAlloc/1024, mem.HeapInuse/1024, mem.HeapIdle/1024, mem.HeapReleased/1024, mem.HeapObjects, a)
			a = 0
		}
	}

	// panic("show me the stacks")

	// On macbook, run run 300 seconds,  private memory is stabilized at 836k
	// On windows, private memory is almost stable.
}
