package collectors

import (
	"fmt"
	"github.com/oliveagle/hickwall/collectors/config"
	"github.com/oliveagle/hickwall/lib/pdh"
	"github.com/oliveagle/hickwall/logging"
	"github.com/oliveagle/hickwall/newcore"
	"strings"
	"time"
)

var (
	_ = fmt.Sprint("")
)

type win_pdh_collector struct {
	name     string // collector name
	interval time.Duration
	enabled  bool
	prefix   string

	// win_pdh_collector specific attributes
	config      config.Config_win_pdh_collector
	hPdh        pdh.PdhCollector
	map_queries map[string]config.Config_win_pdh_query
}

func NewWinPdhCollector(name, prefix string, opts config.Config_win_pdh_collector) newcore.Collector {

	c := win_pdh_collector{
		name:        name,
		enabled:     true,
		prefix:      prefix,
		interval:    opts.Interval.MustDuration(time.Second),
		config:      opts,
		hPdh:        pdh.NewPdhCollector(),
		map_queries: make(map[string]config.Config_win_pdh_query),
	}

	for _, q := range opts.Queries {
		c.hPdh.AddEnglishCounter(q.Query)
		if q.Tags == nil {
			q.Tags = newcore.AddTags.Copy()
		}

		if !q.Ignore_query_tag {
			q.Tags["query"] = q.Query
		}

		c.map_queries[q.Query] = q
	}

	return c
}

func (c win_pdh_collector) Name() string {
	return c.name
}

func (c win_pdh_collector) Close() error {
	c.hPdh.Close()
	return nil
}

func (c win_pdh_collector) ClassName() string {
	return "win_pdh_collector"
}

func (c win_pdh_collector) IsEnabled() bool {
	return c.enabled
}

func (c win_pdh_collector) Interval() time.Duration {
	return c.interval
}

func (c win_pdh_collector) CollectOnce() newcore.CollectResult {
	logging.Debug("win_pdh_collector.CollectOnce Started")

	var items newcore.MultiDataPoint

	for _, pd := range c.hPdh.CollectData() {
		if pd.Err == nil {
			query, ok := c.map_queries[pd.Query]
			if ok == true {
				items = append(items, newcore.NewDP(c.prefix, query.Metric.Clean(), pd.Value, query.Tags, "", "", ""))
			}
		} else {
			if strings.Index(pd.Err.Error(), `\Process(hickwall)\Working Set - Private`) < 0 {
				logging.Errorf("win_pdh_collector ERROR: ", pd.Err)
			}
		}
	}

	logging.Debugf("win_pdh_collector.CollectOnce Finished. count: %d", len(items))
	return newcore.CollectResult{
		Collected: items,
		Next:      time.Now().Add(c.interval),
		Err:       nil,
	}
}
