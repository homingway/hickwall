package collectors

import (
	// "fmt"
	"github.com/oliveagle/hickwall/collectorlib"
	"github.com/oliveagle/hickwall/config"
	"github.com/oliveagle/hickwall/utils"
	log "github.com/oliveagle/seelog"
	"runtime"
	"time"
)

func init() {
	defer utils.Recover_and_log()

	collector_factories["hickwall_client"] = factory_hickwall
}

func factory_hickwall(name string, conf interface{}) <-chan Collector {
	defer utils.Recover_and_log()

	log.Debug("factory_hickwall")

	var out = make(chan Collector)
	go func() {
		// var interval = time.Duration(1) * time.Second
		var interval = time.Duration(1) * time.Millisecond

		out <- &IntervalCollector{
			F:            C_hickwall,
			name:         "factory_hickwall",
			states:       nil,
			Interval:     interval,
			factory_name: "factory_hickwall",
		}
		close(out)
	}()
	return out
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
