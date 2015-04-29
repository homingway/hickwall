package collectors

import (
	"fmt"
	"github.com/oliveagle/hickwall/collectorlib"
	"github.com/oliveagle/hickwall/config"
	log "github.com/oliveagle/seelog"
	"runtime"
	"time"
)

func init() {

	var runtime_conf = config.GetRuntimeConf()

	interval := time.Duration(1) * time.Second
	if runtime_conf.Client_metric_interval != "" {
		ival, err := collectorlib.ParseInterval(runtime_conf.Client_metric_interval)
		if err != nil {
			log.Errorf("cannot parse interval of client_metric_interval: %s - %v", runtime_conf.Client_metric_interval, err)
		}
		interval = ival
	}

	builtin_collectors = append(builtin_collectors, &IntervalCollector{
		F: C_hickwall,
		Enable: func() bool {
			fmt.Println("c_hickwall: enabled: ", runtime_conf.Client_metric_enabled)
			return runtime_conf.Client_metric_enabled
		},
		name:     "builtin_hickwall_client",
		states:   nil,
		Interval: interval,
	})
}

// hickwall process metrics, only runtime stats
func C_hickwall(states interface{}) (collectorlib.MultiDataPoint, error) {
	var (
		md           collectorlib.MultiDataPoint
		m            runtime.MemStats
		runtime_conf = config.GetRuntimeConf()
	)

	tags := AddTags.Copy().Merge(runtime_conf.Tags)
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
