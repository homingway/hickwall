package collectors

import (
	"fmt"
	"github.com/oliveagle/hickwall/newcore"
	"runtime"
	"time"
)

var (
	_ = fmt.Sprint("")

	mem_stats runtime.MemStats
)

type hickwall_collector struct {
	interval time.Duration
	enabled  bool

	// hickwall_collector specific attributes
}

func NewHickwallCollector(interval string) newcore.Collector {
	c := hickwall_collector{
		enabled:  true,
		interval: newcore.NewInterval(interval).MustDuration(time.Second),
	}
	return c
}

func (c hickwall_collector) Name() string {
	return "hickwall_collector"
}

func (c hickwall_collector) Close() error {
	return nil
}

func (c hickwall_collector) ClassName() string {
	return "hickwall_collector"
}

func (c hickwall_collector) IsEnabled() bool {
	return c.enabled
}

func (c hickwall_collector) Interval() time.Duration {
	return c.interval
}

func (c hickwall_collector) CollectOnce() newcore.CollectResult {
	var md newcore.MultiDataPoint

	runtime.ReadMemStats(&mem_stats)

	md = append(md, newcore.NewDP("hickwall.client", "NumGoroutine", runtime.NumGoroutine(), nil, "", "", ""))
	md = append(md, newcore.NewDP("hickwall.client", "mem.Alloc", mem_stats.Alloc, nil, "", "", ""))
	md = append(md, newcore.NewDP("hickwall.client", "mem.TotalAlloc", mem_stats.TotalAlloc, nil, "", "", ""))
	md = append(md, newcore.NewDP("hickwall.client", "mem.Heap.Sys", mem_stats.HeapSys, nil, "", "", ""))
	md = append(md, newcore.NewDP("hickwall.client", "mem.Heap.Alloc", mem_stats.HeapAlloc, nil, "", "", ""))
	md = append(md, newcore.NewDP("hickwall.client", "mem.Heap.Idle", mem_stats.HeapIdle, nil, "", "", ""))
	md = append(md, newcore.NewDP("hickwall.client", "mem.Heap.Inuse", mem_stats.HeapInuse, nil, "", "", ""))
	md = append(md, newcore.NewDP("hickwall.client", "mem.Heap.Released", mem_stats.HeapReleased, nil, "", "", ""))
	md = append(md, newcore.NewDP("hickwall.client", "mem.Heap.Objects", mem_stats.HeapObjects, nil, "", "", ""))
	md = append(md, newcore.NewDP("hickwall.client", "mem.GC.NextGC", mem_stats.NextGC, nil, "", "", ""))
	md = append(md, newcore.NewDP("hickwall.client", "mem.GC.LastGC", mem_stats.LastGC, nil, "", "", ""))
	md = append(md, newcore.NewDP("hickwall.client", "mem.GC.NumGC", mem_stats.NumGC, nil, "", "", ""))
	md = append(md, newcore.NewDP("hickwall.client", "mem.GC.EnableGC", mem_stats.EnableGC, nil, "", "", ""))

	// Add(&md, "hickwall.client", "NumGoroutine", runtime.NumGoroutine(), nil, "", "", "")
	// Add(&md, "hickwall.client", "mem.Alloc", mem_stats.Alloc, nil, "", "", "")
	// Add(&md, "hickwall.client", "mem.TotalAlloc", mem_stats.TotalAlloc, nil, "", "", "")
	// Add(&md, "hickwall.client", "mem.Heap.Sys", mem_stats.HeapSys, nil, "", "", "")
	// Add(&md, "hickwall.client", "mem.Heap.Alloc", mem_stats.HeapAlloc, nil, "", "", "")
	// Add(&md, "hickwall.client", "mem.Heap.Idle", mem_stats.HeapIdle, nil, "", "", "")
	// Add(&md, "hickwall.client", "mem.Heap.Inuse", mem_stats.HeapInuse, nil, "", "", "")
	// Add(&md, "hickwall.client", "mem.Heap.Released", mem_stats.HeapReleased, nil, "", "", "")
	// Add(&md, "hickwall.client", "mem.Heap.Objects", mem_stats.HeapObjects, nil, "", "", "")
	// Add(&md, "hickwall.client", "mem.GC.NextGC", mem_stats.NextGC, nil, "", "", "")
	// Add(&md, "hickwall.client", "mem.GC.LastGC", mem_stats.LastGC, nil, "", "", "")
	// Add(&md, "hickwall.client", "mem.GC.NumGC", mem_stats.NumGC, nil, "", "", "")
	// Add(&md, "hickwall.client", "mem.GC.EnableGC", mem_stats.EnableGC, nil, "", "", "")

	return newcore.CollectResult{
		Collected: md,
		Next:      time.Now().Add(c.interval),
		Err:       nil,
	}
}
