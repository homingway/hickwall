package collectors

// import (
// 	"fmt"
// 	"github.com/oliveagle/hickwall/collectorlib"
// 	"github.com/oliveagle/hickwall/config"
// 	log "github.com/oliveagle/seelog"
// 	"runtime"
// 	"time"
// )

func init() {
	// collector_factories["ping"] = factory_ping
}

// func factory_ping(name string, conf interface{}) <-chan Collector {
// 	log.Debugf("factory_ping, name: %s", name)

// 	var out = make(chan Collector)
// 	go func() {

// 		var (
// 			cf               config.Conf_ping
// 			default_interval = time.Duration(1) * time.Second

// 			runtime_conf = config.GetRuntimeConf()
// 		)

// 		if conf != nil {
// 			interval, err := collectorlib.ParseInterval(cf.Interval)
// 			if err != nil {
// 				log.Errorf("cannot parse interval of collector_pdh: %s - %v", cf.Interval, err)
// 				interval = default_interval
// 			}

// 			for idx, target := range cf.Targets {

// 				var (
// 					states state_c_ping
// 				)

// 				states.Interval = interval
// 				states.Target = target

// 				out <- &IntervalCollector{
// 					F:        C_ping,
// 					Enable:   nil,
// 					name:     fmt.Sprintf("%s_%d", name, idx),
// 					states:   states,
// 					Interval: states.Interval,
// 				}
// 			}
// 		}

// 		close(out)
// 	}()

// 	return out
// }

// type state_c_ping struct {
// 	Interval time.Duration
// 	Target   string
// }

// // hickwall process metrics, only runtime stats
// func C_ping(states interface{}) (collectorlib.MultiDataPoint, error) {
// 	var (
// 		md           collectorlib.MultiDataPoint
// 		m            runtime.MemStats
// 		runtime_conf = config.GetRuntimeConf()
// 	)

// 	tags := AddTags.Copy().Merge(runtime_conf.Tags)
// 	// runtime.ReadMemStats(&m)

// 	// Add(&md, "hickwall.client.NumGoroutine", runtime.NumGoroutine(), tags, "", "", "")
// 	return md, nil
// }
