package main

import (
	"fmt"
	// "runtime/pprof"
	// "github.com/kr/pretty"
	"runtime"
	"time"
)

func main() {
	done := time.After(time.Second * 60)
	tick := time.Tick(time.Second * 1)
	m := runtime.MemStats{}
loop:
	for {
		select {
		case <-tick:
			fmt.Println("----------------------------------------------------------------------")
			// fmt.Println("NumCgoCall: ", runtime.NumCgoCall())
			// fmt.Println("NumGoroutine: ", runtime.NumGoroutine())
			// runtime.ReadMemStats(&m)
			// fmt.Println(m)
			// pretty.Println(m)

			// fmt.Println("MemStats.Sys", m.Sys/1024) // bytes obtained from system (sum of XxxSys below)

			// fmt.Println("MemStats.Alloc:  kb", m.Alloc/1024)           // bytes allocated and still in use
			// fmt.Println("MemStats.TotalAlloc:  kb", m.TotalAlloc/1024) // bytes allocated (even if freed)

			// fmt.Println("MemStats.HeapSys  kb", m.HeapSys/1024)           // bytes obtained from system
			// fmt.Println("MemStats.HeapAlloc  kb", m.HeapAlloc/1024)       // bytes allocated and still in use      !!!
			// fmt.Println("MemStats.HeapIdle  kb", m.HeapIdle/1024)         // bytes in idle spans
			// fmt.Println("MemStats.HeapInuse  kb", m.HeapInuse/1024)       // bytes in non-idle span
			// fmt.Println("MemStats.HeapReleased  kb", m.HeapReleased/1024) // bytes released to the OS
			// fmt.Println("MemStats.HeapObjects", m.HeapObjects) // total number of allocated objects

			// fmt.Println("MemStats.NextGC  kb", m.NextGC/1024) // next collection will happen when HeapAlloc ≥ this amount
			// fmt.Println("MemStats.LastGC  ns", m.LastGC)      // end time of last collection (nanoseconds since 1970)

			// Main allocation heap statistics.

			// // Garbage collector statistics.
			// fmt.Println("MemStats.NextGC", m.NextGC) // next collection will happen when HeapAlloc ≥ this amount
			// fmt.Println("MemStats.LastGC", m.LastGC) // end time of last collection (nanoseconds since 1970)

			// fmt.Println("MemStats.NumGC", m.NumGC)
			// fmt.Println("MemStats.EnableGC", m.EnableGC)
			// fmt.Println("MemStats.DebugGC", m.DebugGC)

			// // Per-size allocation statistics.
			// // 61 is NumSizeClasses in the C code.
			// BySize [61]struct {
			//         Size    uint32
			//         Mallocs uint64
			//         Frees   uint64
			// }

		case <-done:
			break loop
		}
	}
}
