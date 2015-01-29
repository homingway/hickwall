package main

import (
	// log "github.com/oliveagle/hickwall/_third_party/seelog"
	// "github.com/oliveagle/hickwall/config"
	"github.com/rcrowley/go-metrics"
	"testing"
	"time"
)

func TestLogMetric(t *testing.T) {

	c := metrics.NewCounter()
	metrics.Register("foo", c)
	c.Inc(12)
	go func() {
		for {
			c.Inc(1)
			time.Sleep(10 * time.Millisecond)
		}
	}()

	go func() {
		for {
			t.Log(c.Count())
			time.Sleep(1 * time.Second)
		}
	}()
	time.Sleep(3 * time.Second)

	t.Error("---------")

}
