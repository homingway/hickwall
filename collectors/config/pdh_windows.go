package config

import (
	"github.com/oliveagle/hickwall/newcore"
)

type Config_win_pdh_query struct {
	Query        string         `json:"query"`
	Metric       newcore.Metric `json:"metric"`
	Tags         newcore.TagSet `json:"tags"`
	Query_to_tag bool           `json:"query_to_tag"` // true, query will be add to tags.
}

type Config_win_pdh_collector struct {
	Interval     newcore.Interval       `json:"interval"`
	Tags         newcore.TagSet         `json:"tags"`
	Queries      []Config_win_pdh_query `json:"queries"`
	Query_to_tag bool                   `json:"query_to_tag"` // true, query will be add to tags.
}
