package collectors

import (
	"fmt"
	"github.com/oliveagle/hickwall/newcore"
	"runtime"
	"time"
)

var (
	_ = fmt.Sprint("")
)

type hickwall_collector struct {
	interval time.Duration
	enabled  bool

	// hickwall_collector specific attributes
	mem_stats runtime.MemStats
}

func NewHickwallCollector(interval string) newcore.Collector {
	c := &hickwall_collector{
		enabled:  true,
		interval: newcore.NewInterval(interval).MustDuration(time.Second),
	}
	return c
}

func (c *hickwall_collector) Name() string {
	return "hickwall_collector"
}

func (c *hickwall_collector) Close() error {
	return nil
}

func (c *hickwall_collector) ClassName() string {
	return "hickwall_collector"
}

func (c *hickwall_collector) IsEnabled() bool {
	return c.enabled
}

func (c *hickwall_collector) Interval() time.Duration {
	return c.interval
}

func (c *hickwall_collector) CollectOnce() newcore.CollectResult {
	var md newcore.MultiDataPoint

	runtime.ReadMemStats(&c.mem_stats)

	Add(&md, "hickwall.client", "NumGoroutine", runtime.NumGoroutine(), nil, "", "", "")
	Add(&md, "hickwall.client", "mem.Alloc", c.mem_stats.Alloc, nil, "", "", "")
	Add(&md, "hickwall.client", "mem.TotalAlloc", c.mem_stats.TotalAlloc, nil, "", "", "")
	Add(&md, "hickwall.client", "mem.Heap.Sys", c.mem_stats.HeapSys, nil, "", "", "")
	Add(&md, "hickwall.client", "mem.Heap.Alloc", c.mem_stats.HeapAlloc, nil, "", "", "")
	Add(&md, "hickwall.client", "mem.Heap.Idle", c.mem_stats.HeapIdle, nil, "", "", "")
	Add(&md, "hickwall.client", "mem.Heap.Inuse", c.mem_stats.HeapInuse, nil, "", "", "")
	Add(&md, "hickwall.client", "mem.Heap.Released", c.mem_stats.HeapReleased, nil, "", "", "")
	Add(&md, "hickwall.client", "mem.Heap.Objects", c.mem_stats.HeapObjects, nil, "", "", "")
	Add(&md, "hickwall.client", "mem.GC.NextGC", c.mem_stats.NextGC, nil, "", "", "")
	Add(&md, "hickwall.client", "mem.GC.LastGC", c.mem_stats.LastGC, nil, "", "", "")
	Add(&md, "hickwall.client", "mem.GC.NumGC", c.mem_stats.NumGC, nil, "", "", "")
	Add(&md, "hickwall.client", "mem.GC.EnableGC", c.mem_stats.EnableGC, nil, "", "", "")

	return newcore.CollectResult{
		Collected: md,
		Next:      time.Now().Add(c.interval),
		Err:       nil,
	}
}
