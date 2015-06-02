package config

import (
	c_conf "github.com/oliveagle/hickwall/collectors/config"
)

type CollectorConfigGroup struct {
	Prefix         string               `json:"prefix"`
	Collector_ping []c_conf.Config_Ping `json:"collector_ping"`
	//  Collector_cmd     []c_conf.Config_command           `json:"collector_cmd"`

	//  Collector_win_sys Conf_win_sys `json:"collector_win_sys"`
}
