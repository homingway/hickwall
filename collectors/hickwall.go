package collectors

import (
	"github.com/oliveagle/hickwall/collectorlib"
	"github.com/oliveagle/hickwall/config"
	"runtime"
	// "time"
)

// func init() {

// 	interval := time.Duration(1) * time.Second
// 	if config.Conf.Client_metric_interval > 0 {
// 		interval = time.Duration(config.Conf.Client_metric_interval) * time.Second
// 	}
// 	// fmt.Println("interval :", interval)

// 	builtin_collectors = append(builtin_collectors, &IntervalCollector{
// 		F: c_hickwall,
// 		Enable: func() bool {
// 			// fmt.Println("config.Conf.Client_metric_enabled: ", config.Conf.Client_metric_enabled)
// 			return config.Conf.Client_metric_enabled
// 		},
// 		name:     "builtin_hickwall_client",
// 		states:   nil,
// 		Interval: interval,
// 	})
// }

// hickwall process metrics, only runtime stats
func C_hickwall(states interface{}) (collectorlib.MultiDataPoint, error) {
	var md collectorlib.MultiDataPoint
	var m runtime.MemStats

	tags := AddTags.Copy().Merge(config.Conf.Tags)
	runtime.ReadMemStats(&m)

	Add(&md, "hickwall.client.NumGoroutine", runtime.NumGoroutine(), tags, "", "", "")
	Add(&md, "hickwall.client.mem.Alloc", m.Alloc, tags, "", "", "")
	Add(&md, "hickwall.client.mem.TotalAlloc", m.TotalAlloc, tags, "", "", "")
	Add(&md, "hickwall.client.mem.Heap.Sys", m.HeapSys, tags, "", "", "")
	Add(&md, "hickwall.client.mem.Heap.Alloc", m.HeapAlloc, tags, "", "", "")
	Add(&md, "hickwall.client.mem.Heap.Idle", m.HeapIdle, tags, "", "", "")
	Add(&md, "hickwall.client.mem.Heap.Inuse", m.HeapInuse, tags, "", "", "")
	Add(&md, "hickwall.client.mem.Heap.Released", m.HeapReleased, tags, "", "", "")
	Add(&md, "hickwall.client.mem.Heap.Objects", m.HeapObjects, tags, "", "", "")
	Add(&md, "hickwall.client.mem.GC.NextGC", m.NextGC, tags, "", "", "")
	Add(&md, "hickwall.client.mem.GC.LastGC", m.LastGC, tags, "", "", "")
	Add(&md, "hickwall.client.mem.GC.NumGC", m.NumGC, tags, "", "", "")
	Add(&md, "hickwall.client.mem.GC.EnableGC", m.EnableGC, tags, "", "", "")
	return md, nil
}
