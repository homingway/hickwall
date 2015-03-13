package collectorlib

import (
	"encoding/json"
	// "fmt"
	"math"
	"math/big"
	// "strconv"
	"time"
)

//NOTE: these struct don't have any use case at the time.

var bigMaxInt64 = big.NewInt(math.MaxInt64)

type DataPoint struct {
	Metric    string      `json:"metric"`
	Timestamp time.Time   `json:"timestamp"`
	Value     interface{} `json:"value"`
	Tags      TagSet      `json:"tags"`
}

func (d *DataPoint) MarshalJSON() ([]byte, error) {

	return json.Marshal(struct {
		Metric    string      `json:"metric"`
		Timestamp time.Time   `json:"timestamp"`
		Value     interface{} `json:"value"`
		Tags      TagSet      `json:"tags"`
	}{
		d.Metric,
		d.Timestamp,
		d.Value,
		d.Tags,
	})
}

func (d *DataPoint) GetFlatMetric(tpl string) (string, error) {
	return FlatMetricKeyAndTags(tpl, d.Metric, d.Tags)
}

type TagSet map[string]string

// func (t *TagSet) Copy(t1 *TagSet) *TagSet {

// }
