package config

import (
	"github.com/oliveagle/hickwall/newcore"
)

type Config_win_wmi struct {
	Tags     newcore.TagSet         `json:"tags"`
	Interval newcore.Interval       `json:"interval"`
	Queries  []Config_win_wmi_query `json:"queries"`
}
type Config_win_wmi_query struct {
	Query   string                        `json:"query"`
	Tags    newcore.TagSet                `json:"tags"`
	Metrics []Config_win_wmi_query_metric `json:"metrics"`
}
type Config_win_wmi_query_metric struct {
	//TODO: Meta
	Value_from string            `json:"value_from"`
	Metric     newcore.Metric    `json:"metric"`
	Tags       newcore.TagSet    `json:"tags"`
	Meta       map[string]string `json:"meta"`
	Default    interface{}       `json:"default"`
}
