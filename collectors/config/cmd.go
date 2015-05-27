package config

import (
	"github.com/oliveagle/hickwall/newcore"
)

type Config_command struct {
	Metric   newcore.Metric   `json:"metric"`
	Cmd      []string         `json:"cmd"`
	Interval newcore.Interval `json:"interval"` // default 1s
	Tags     newcore.TagSet   `json:"tags"`
	Timeout  newcore.Interval `json:"timeout"`
}
