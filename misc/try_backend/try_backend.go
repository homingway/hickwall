package main

import (
	"github.com/oliveagle/hickwall/backends"
	"github.com/oliveagle/hickwall/collectorlib"

	"fmt"
	"math/rand"
	"time"
)

func main() {

	fmt.Println("--")
	// p := collectorlib.DataPoint{
	// 	Metric:    "metric1",
	// 	Timestamp: time.Now().UnixNano(),
	// 	Value:     1,
	// }
	// md := collectorlib.MultiDataPoint{&p}
	// fmt.Println(md)

	// backend, _ := backends.GetBackendByName("stdout")
	// go backend.Run()
	// defer backend.Close()

	fmt.Println(backends.GetBackendList())

	backends.RunBackends()
	defer backends.CloseBackends()

	tick := time.Tick(time.Millisecond * 10)
	done := time.After(time.Second * 60)

loop:
	for {
		select {
		case <-tick:
			// backend.Write(md)
			rand.Seed(time.Now().UnixNano())
			p := collectorlib.DataPoint{
				Metric:    "metric1",
				Timestamp: time.Now(),
				Value:     rand.Intn(100),
			}
			md := collectorlib.MultiDataPoint{&p}
			backends.WriteToBackends(md)
		case <-done:
			break loop
		}
	}

}
