package newcore

import (
	"time"
)

type DataPoint struct {
	Metric    string            `json:"metric"`
	Timestamp time.Time         `json:"timestamp"`
	Value     interface{}       `json:"value"`
	Tags      map[string]string `json:"tags"`
	Meta      map[string]string `json:"meta,omitempty"`
}

type MultiDataPoint []*DataPoint
