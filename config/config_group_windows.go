package config

import (
	c_conf "github.com/oliveagle/hickwall/collectors/config"
)

type CollectorConfigGroup struct {
	Prefix            string                            `json:"prefix"`
	Collector_ping    []c_conf.Config_Ping              `json:"collector_ping"`
	Collector_win_pdh []c_conf.Config_win_pdh_collector `json:"collector_win_pdh"`
	Collector_win_wmi []c_conf.Config_win_wmi           `json:"collector_win_wmi"`
	Collector_win_sys *c_conf.Config_win_sys            `json:"collector_win_sys"`
}
