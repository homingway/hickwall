package main

import (
	"github.com/oliveagle/go-collectors/datapoint"
	"github.com/oliveagle/hickwall/backends"

	"fmt"
	"time"
)

func main() {

	p := datapoint.DataPoint{
		Metric:    "metric1",
		Timestamp: time.Now().UnixNano(),
		Value:     1,
	}
	md := datapoint.MultiDataPoint{&p}
	fmt.Println(md)

	// backend, _ := backends.GetBackendByName("stdout")
	// go backend.Run()
	// defer backend.Close()

	backends.RunBackends()
	defer backends.CloseBackends()

	tick := time.Tick(time.Millisecond * 10)
	done := time.After(time.Second * 5)

loop:
	for {
		select {
		case <-tick:
			// backend.Write(md)
			backends.WriteToBackends(md)
		case <-done:
			break loop
		}
	}

}
