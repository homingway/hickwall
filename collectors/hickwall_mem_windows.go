package collectors

import (
	"github.com/oliveagle/hickwall/collectors/config"
	"github.com/oliveagle/hickwall/newcore"
)

func NewWinHickwallMemCollector(interval string, tags newcore.TagSet) newcore.Collector {
	opts := config.Config_win_pdh_collector{
		Interval: newcore.Interval(interval),
		Tags:     tags,
		Queries: []config.Config_win_pdh_query{
			{
				Query:            "\\Process(hickwall)\\Working Set - Private",
				Metric:           "private_working_set.bytes",
				Ignore_query_tag: true,
			},
		},
	}
	return NewWinPdhCollector("hickwall_mem", "hickwall.client.mem", opts)
}
