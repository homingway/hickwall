package main

// import (
// 	// "fmt"
// 	log "github.com/cihub/seelog"
// 	// "github.com/oliveagle/hickwall/config"
// 	// "github.com/oliveagle/go-metrics"
// 	"net"
// 	"time"
// )

// // Output each metric in the given registry periodically using the given
// // logger.
// func Log(r metrics.Registry, d time.Duration) {
// 	for _ = range time.Tick(d) {
// 		collectors.PrintMetrics(r, log.Info)
// 	}
// }

// func registry_metrics() {
// 	c := metrics.NewCounter()
// 	metrics.Register("foo", c)
// }

// func serve_graphite() {
// 	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:2003")
// 	go metrics.Graphite(metrics.DefaultRegistry, 1*time.Second, "metrics", addr)
// }

// func process_tick(tick chan time.Time) {
// 	log.Info("process_tick")
// 	// metrics.DefaultRegistry.GetOrRegister()

// 	// go func() {
// 	// 	for {
// 	// 		c.Inc(47)
// 	// 		time.Sleep(10 * time.Millisecond)
// 	// 	}
// 	// }()

// 	// go Log(metrics.DefaultRegistry, 1*time.Second)
// }
