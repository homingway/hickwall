// +build windows

package collectors

import (
	// "fmt"
	"github.com/oliveagle/go-collectors/datapoint"
	"github.com/oliveagle/go-collectors/pdh"
	// "github.com/oliveagle/hickwall/collectorlib"
	"github.com/oliveagle/hickwall/config"
	// "github.com/kr/pretty"
	// log "github.com/cihub/seelog"
	"time"
)

func init() {
	collector_factories["win_pdh"] = factory_win_pdh

	builtin_collectors = append(builtin_collectors, builtin_win_pdh())
}

func builtin_win_pdh() Collector {
	gauge := pdh.NewPdhCollector()
	key_map := make(map[string]string)
	//TODO: defer gauge.Close()

	gauge.AddEnglishCounter("\\System\\Processes")
	key_map["\\System\\Processes"], _ = ValidateKey("win.processes.count")

	gauge.AddEnglishCounter("\\Memory\\Available Bytes")
	key_map["\\Memory\\Available Bytes"], _ = ValidateKey("win.memory.available_bytes")

	gauge.AddEnglishCounter("\\Processes(_Total)\\Working Set")
	key_map["\\Processes(_Total)\\Working Set"], _ = ValidateKey("win.processes.working_set.total")

	gauge.AddEnglishCounter("\\Memory\\Cache Bytes")
	key_map["\\Memory\\Cache Bytes"], _ = ValidateKey("win.memory.cache_bytes")

	defaultStates := MockStates{
		Interval:  time.Second * 1,
		MetricKey: "os",

		hPdh:    gauge,
		key_map: key_map,
	}

	return &IntervalCollector{
		F:        c_win_pdh,
		Enable:   nil,
		name:     "builtin_win_pdh",
		states:   defaultStates,
		Interval: defaultStates.Interval,
	}
}

func factory_win_pdh(name string, conf interface{}) Collector {
	// Conf_win_pdh
	var states MockStates
	var cf config.Conf_win_pdh

	if conf != nil {
		cf = conf.(config.Conf_win_pdh)
		// fmt.Println("factory_win_pdh: ", cf)

		states.hPdh = pdh.NewPdhCollector()
		states.key_map = make(map[string]string)
		states.Interval = time.Duration(cf.Interval) * time.Second
		states.Tags = datapoint.TagSet{}

		for query, metric_key := range cf.Queries {
			// fmt.Println(query, metric_key)
			states.hPdh.AddEnglishCounter(query)
			// states.key_map[query] = metric_key
			states.key_map[query], _ = ValidateKey(metric_key)
		}

		for key, value := range config.Conf.Tags {
			states.Tags[key] = value
		}

		// same key will override
		for key, value := range cf.Tags {
			states.Tags[key] = value
		}
	}

	return &IntervalCollector{
		F:        c_win_pdh,
		Enable:   nil,
		name:     name,
		states:   states,
		Interval: states.Interval,
	}
}

type MockStates struct {
	Interval  time.Duration
	MetricKey string
	Tags      datapoint.TagSet

	// internal use only
	hPdh    *pdh.PdhCollector
	key_map map[string]string
}

func c_win_pdh(states interface{}) (datapoint.MultiDataPoint, error) {
	var md datapoint.MultiDataPoint
	var st MockStates

	if states != nil {
		st = states.(MockStates)
		// fmt.Println("c_win_pdh states: ", states)
	}

	if st.hPdh != nil {

		data := st.hPdh.CollectData()
		key_map := st.key_map

		for _, pd := range data {
			Add(&md, key_map[pd.Query], pd.Value, st.Tags, "", "", "")
		}
	}

	// Add(&md, "this.is.metric.key.string", " # the string value # ", st.Tags, "", "", "")
	// Add(&md, "this.is.metric.key.int", 1, st.Tags, "", "", "")

	return md, nil
}
