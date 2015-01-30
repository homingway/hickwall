package main

import (
	// "fmt"
	log "github.com/cihub/seelog"
	// "github.com/oliveagle/hickwall/config"
	"github.com/rcrowley/go-metrics"
	"net"
	"time"
)

// Output each metric in the given registry periodically using the given
// logger.
func Log(r metrics.Registry, d time.Duration) {
	for _ = range time.Tick(d) {
		r.Each(func(name string, i interface{}) {
			switch metric := i.(type) {
			case metrics.Counter:
				log.Infof("counter %s :    %9d", name, metric.Count())
			case metrics.Gauge:
				log.Infof("gauge %s :    %9d", name, metric.Value())
			case metrics.GaugeFloat64:
				log.Infof("gauge %s :   %f", name, metric.Value())
			case metrics.Healthcheck:
				metric.Check()
				log.Infof("healthcheck %s :     error: %v", name, metric.Error())
			case metrics.Histogram:
				h := metric.Snapshot()
				ps := h.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
				log.Infof("histogram %s", name)
				log.Infof("  count:       %9d", h.Count())
				log.Infof("  min:         %9d", h.Min())
				log.Infof("  max:         %9d", h.Max())
				log.Infof("  mean:        %12.2f", h.Mean())
				log.Infof("  stddev:      %12.2f", h.StdDev())
				log.Infof("  median:      %12.2f", ps[0])
				log.Infof("  75%%:         %12.2f", ps[1])
				log.Infof("  95%%:         %12.2f", ps[2])
				log.Infof("  99%%:         %12.2f", ps[3])
				log.Infof("  99.9%%:       %12.2f", ps[4])
			case metrics.Meter:
				m := metric.Snapshot()
				log.Infof("meter %s", name)
				log.Infof("  count:       %9d", m.Count())
				log.Infof("  1-min rate:  %12.2f", m.Rate1())
				log.Infof("  5-min rate:  %12.2f", m.Rate5())
				log.Infof("  15-min rate: %12.2f", m.Rate15())
				log.Infof("  mean rate:   %12.2f", m.RateMean())
			case metrics.Timer:
				t := metric.Snapshot()
				ps := t.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
				log.Infof("timer %s", name)
				log.Infof("  count:       %9d", t.Count())
				log.Infof("  min:         %9d", t.Min())
				log.Infof("  max:         %9d", t.Max())
				log.Infof("  mean:        %12.2f", t.Mean())
				log.Infof("  stddev:      %12.2f", t.StdDev())
				log.Infof("  median:      %12.2f", ps[0])
				log.Infof("  75%%:         %12.2f", ps[1])
				log.Infof("  95%%:         %12.2f", ps[2])
				log.Infof("  99%%:         %12.2f", ps[3])
				log.Infof("  99.9%%:       %12.2f", ps[4])
				log.Infof("  1-min rate:  %12.2f", t.Rate1())
				log.Infof("  5-min rate:  %12.2f", t.Rate5())
				log.Infof("  15-min rate: %12.2f", t.Rate15())
				log.Infof("  mean rate:   %12.2f", t.RateMean())
			}
		})
	}
}

func registry_metrics() {
	c := metrics.NewCounter()
	metrics.Register("foo", c)
}

func serve_graphite() {
	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:2003")
	go metrics.Graphite(metrics.DefaultRegistry, 1*time.Second, "metrics", addr)
}

func process_tick(tick chan time.Time) {
	log.Info("process_tick")
	// metrics.DefaultRegistry.GetOrRegister()

	// go func() {
	// 	for {
	// 		c.Inc(47)
	// 		time.Sleep(10 * time.Millisecond)
	// 	}
	// }()

	// go Log(metrics.DefaultRegistry, 1*time.Second)
}
