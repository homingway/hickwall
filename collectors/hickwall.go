package collectors

import (
	"fmt"
	"github.com/oliveagle/hickwall/collectorlib"
	"github.com/oliveagle/hickwall/config"
	"github.com/oliveagle/hickwall/utils"
	log "github.com/oliveagle/seelog"
	"runtime"
	"time"
)

func init() {
	defer utils.Recover_and_log()

	client_conf := config.GetRuntimeConf().Client

	interval := time.Duration(1) * time.Second
	if client_conf.Metric_interval != "" {
		ival, err := collectorlib.ParseInterval(client_conf.Metric_interval)
		if err != nil {
			log.Errorf("cannot parse interval of client.Metric_interval: %s - %v", client_conf.Metric_interval, err)
		}
		interval = ival
	}

	builtin_collectors = append(builtin_collectors, &IntervalCollector{
		F: C_hickwall,
		Enable: func() bool {
			fmt.Println("c_hickwall: enabled: ", client_conf.Metric_enabled)
			return client_conf.Metric_enabled
		},
		name:         "hickwall_client",
		states:       nil,
		Interval:     interval,
		factory_name: "hickwall_client",
	})
}

// hickwall process metrics, only runtime stats
func C_hickwall(states interface{}) (collectorlib.MultiDataPoint, error) {
	defer utils.Recover_and_log()

	var (
		md collectorlib.MultiDataPoint
		m  runtime.MemStats
	)

	client_conf := config.GetRuntimeConf().Client

	tags := AddTags.Copy().Merge(client_conf.Tags)
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
