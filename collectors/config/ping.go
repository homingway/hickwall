package config

import (
	"github.com/oliveagle/hickwall/newcore"
)

type Config_single_pinger struct {
	Metric   string           `json:"metric"`
	Tags     newcore.TagSet   `json:"tags"`
	Target   string           `json:"target"`
	Packets  int              `json:"packets"`
	Timeout  newcore.Interval `json:"timeout"`
	Interval newcore.Interval `json:"interval"`
}
