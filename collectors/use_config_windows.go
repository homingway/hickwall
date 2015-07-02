// +build windows
package collectors

import (
	"fmt"
	"github.com/oliveagle/hickwall/collectors/windows"
	"github.com/oliveagle/hickwall/config"
	"github.com/oliveagle/hickwall/logging"
	"github.com/oliveagle/hickwall/newcore"
)

func gen_collector_name(gid, cid int, cname string) string {
	return fmt.Sprintf("g_%d_%s_c_%d", gid, cname, cid)
}

func UseConfigCreateCollectors(rconf *config.RuntimeConfig) ([]newcore.Collector, error) {
	var clrs []newcore.Collector
	var prefixs = make(map[string]bool)

	for gid, group := range rconf.Groups {
		logging.Infof("gid: %d, prefix: %s\n", gid, group.Prefix)
		if len(group.Prefix) <= 0 {
			return nil, fmt.Errorf("group (idx:%d) prefix is empty.", gid)
		} else {

			_, exists := prefixs[group.Prefix]
			if exists == false {
				prefixs[group.Prefix] = true
			} else {
				return nil, fmt.Errorf("duplicated group prefix: %s", group.Prefix)
			}
		}

		for cid, conf := range group.Collector_ping {
			pings := MustNewPingCollectors(gen_collector_name(gid, cid, "ping"), group.Prefix, conf)
			for _, c := range pings {
				clrs = append(clrs, c)
			}
		}

		for cid, conf := range group.Collector_win_pdh {
			c := windows.MustNewWinPdhCollector(gen_collector_name(gid, cid, "pdh"), group.Prefix, conf)
			clrs = append(clrs, c)
		}

		for cid, conf := range group.Collector_win_wmi {
			c := windows.MustNewWinWmiCollector(gen_collector_name(gid, cid, "wmi"), group.Prefix, conf)
			clrs = append(clrs, c)
		}

		if group.Collector_win_sys != nil {
			cs := windows.MustNewWinSysCollectors(gen_collector_name(gid, 0, "win_sys"), group.Prefix, group.Collector_win_sys)
			for _, c := range cs {
				clrs = append(clrs, c)
			}
			//			clrs = append(clrs, cs...)
		}

	}

	logging.Debugf("rconf.Client.Metric_Enabled: %v, rconf.Client.Metric_Interval: %v",
		rconf.Client.Metric_Enabled, rconf.Client.Metric_Interval)
	if rconf.Client.Metric_Enabled == true {
		clrs = append(clrs, MustNewHickwallCollector(rconf.Client.Metric_Interval))
		clrs = append(clrs, windows.MustNewWinHickwallMemCollector(rconf.Client.Metric_Interval, rconf.Client.Tags))
	}
	return clrs[:], nil
}
