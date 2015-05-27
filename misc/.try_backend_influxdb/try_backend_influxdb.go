package main

// import (
// 	"github.com/oliveagle/hickwall/backends"
// 	"github.com/oliveagle/hickwall/collectorlib"
// 	"github.com/oliveagle/hickwall/config"

// 	"fmt"
// 	"math/rand"
// 	"time"
// )

// func main() {
// 	fmt.Println("---")

// 	config.LoadRuntimeConfFromFileOnce()
// 	backends.CreateBackendsFromRuntimeConf()

// 	// backend, _ := backends.GetBackendByName("influxdb")
// 	// backend, _ := backends.GetBackendByNameVersion("influxdb", "0.9.0-rc7")
// 	backend, _ := backends.GetBackendByNameVersion("influxdb", "0.8.8")
// 	go backend.Run()
// 	defer backend.Close()

// 	// backends.RunBackends()
// 	// defer backends.CloseBackends()

// 	tick := time.Tick(time.Millisecond * 1000)
// 	done := time.Tick(time.Second * 6)

// loop:
// 	for {
// 		select {
// 		case <-tick:
// 			fmt.Println(" <- tick ----------------------")
// 			for i := 0; i < 10; i++ {
// 				rand.Seed(time.Now().UTC().UnixNano())
// 				p := collectorlib.DataPoint{
// 					Metric:    fmt.Sprintf("metric1.%d", i),
// 					Timestamp: time.Now(),
// 					Value:     rand.Float64(),
// 					Tags: map[string]string{
// 						"bu": "hotel",
// 					},
// 				}
// 				md := collectorlib.MultiDataPoint{p}
// 				backend.Write(md)
// 			}
// 			// backends.WriteToBackends(md)
// 		case <-done:
// 			fmt.Println(" <- done --------------------- done -")
// 			break loop
// 		}
// 	}
// }
