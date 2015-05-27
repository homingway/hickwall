package collectors

import (
	"fmt"
	"github.com/oliveagle/hickwall/config"
	"github.com/oliveagle/hickwall/newcore"
)

func gen_collector_name(gid, cid int, cname string) string {
	return fmt.Sprintf("g_%d_%s_c_%d", gid, cname, cid)
}

func UseConfigCreateCollectors(rconf *config.RuntimeConfig) ([]newcore.Collector, error) {
	var clrs []newcore.Collector
	var prefixs = make(map[string]bool)

	if rconf != nil {
		for gid, group := range rconf.Groups {

			fmt.Printf("gid: %d, prefix: %s", gid, group.Prefix)
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
				c := NewPingCollectors(gen_collector_name(gid, cid, "ping"), group.Prefix, conf)
				clrs = append(clrs, c...)
			}

			for cid, conf := range group.Collector_cmd {
				c := NewCmdCollector(gen_collector_name(gid, cid, "cmd"), group.Prefix, conf)
				clrs = append(clrs, c)
			}

			for cid, conf := range group.Collector_win_pdh {
				c := NewWinPdhCollector(gen_collector_name(gid, cid, "pdh"), group.Prefix, conf)
				clrs = append(clrs, c)
			}

			for cid, conf := range group.Collector_win_wmi {
				c := NewWinWmiCollector(gen_collector_name(gid, cid, "wmi"), group.Prefix, conf)
				clrs = append(clrs, c)
			}
		}

		return clrs[:], nil
	} else {
		return nil, fmt.Errorf("rconf is nil")
	}
}
