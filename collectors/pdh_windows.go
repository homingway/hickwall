// +build windows

package collectors

import (
	// "fmt"
	"github.com/oliveagle/go-collectors/datapoint"
	"github.com/oliveagle/go-collectors/pdh"
	// "github.com/oliveagle/hickwall/collectorlib"
	// "github.com/kr/pretty"
	"github.com/oliveagle/hickwall/config"
	// log "github.com/cihub/seelog"
	"time"
)

func init() {
	collector_factories["win_pdh"] = factory_win_pdh

	builtin_collectors = append(builtin_collectors, builtin_win_pdh())
}

func builtin_win_pdh() Collector {

	queries := []config.Conf_win_pdh_query{}

	queries = append(queries, config.Conf_win_pdh_query{
		Query:  "\\System\\Processes",
		Metric: "win.processes.count"})
	queries = append(queries, config.Conf_win_pdh_query{
		Query:  "\\Memory\\Available Bytes",
		Metric: "win.memory.available_bytes"})
	queries = append(queries, config.Conf_win_pdh_query{
		Query:  "\\Processes(_Total)\\Working Set",
		Metric: "win.processes.working_set.total"})
	queries = append(queries, config.Conf_win_pdh_query{
		Query:  "\\Memory\\Cache Bytes",
		Metric: "win.memory.cache_bytes"})

	conf := config.Conf_win_pdh{
		Interval: 2,
		Queries:  queries,
	}

	return factory_win_pdh("builtin_win_pdh", conf)
}

func factory_win_pdh(name string, conf interface{}) Collector {
	var states state_win_pdh
	var cf config.Conf_win_pdh

	if conf != nil {
		cf = conf.(config.Conf_win_pdh)
		// fmt.Println("factory_win_pdh: ", cf)
		// pretty.Println("factory_win_pdh:", cf)

		states.hPdh = pdh.NewPdhCollector()
		states.Interval = time.Duration(cf.Interval) * time.Second
		states.map_queries = make(map[string]config.Conf_win_pdh_query)

		for _, query_obj := range cf.Queries {
			query := query_obj.Query
			//TODO: validate query

			states.hPdh.AddEnglishCounter(query)

			query_obj.Tags = AddTags.Copy().Merge(config.Conf.Tags).Merge(cf.Tags).Merge(query_obj.Tags)

			states.map_queries[query] = query_obj
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

type state_win_pdh struct {
	Interval time.Duration

	// internal use only
	hPdh        *pdh.PdhCollector
	map_queries map[string]config.Conf_win_pdh_query
}

func c_win_pdh(states interface{}) (datapoint.MultiDataPoint, error) {
	var md datapoint.MultiDataPoint
	var st state_win_pdh

	if states != nil {
		st = states.(state_win_pdh)
		// fmt.Println("c_win_pdh states: ", states)
	}

	if st.hPdh != nil {

		data := st.hPdh.CollectData()
		queries := st.map_queries

		for _, pd := range data {
			query := queries[pd.Query]
			Add(&md, query.Metric, pd.Value, query.Tags, "", "", "")
		}
	}

	// Add(&md, "this.is.metric.key.string", " # the string value # ", nil, "", "", "")
	// Add(&md, "this.is.metric.key.int", 1, nil, "", "", "")

	return md, nil
}
