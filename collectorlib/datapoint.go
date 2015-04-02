package collectorlib

import (
	"encoding/json"
	// "fmt"
	"math"
	"math/big"
	"strconv"
	"time"
)

type MultiDataPoint []DataPoint

var bigMaxInt64 = big.NewInt(math.MaxInt64)

type DataPoint struct {
	Metric    string            `json:"metric"`
	Timestamp time.Time         `json:"timestamp"`
	Value     interface{}       `json:"value"`
	Tags      TagSet            `json:"tags"`
	Meta      map[string]string `json:"meta"`
}

func (d *DataPoint) MarshalJSON() ([]byte, error) {

	d.Clean()

	return json.Marshal(struct {
		Metric    string            `json:"metric"`
		Timestamp time.Time         `json:"timestamp"`
		Value     interface{}       `json:"value"`
		Tags      TagSet            `json:"tags"`
		Meta      map[string]string `json:"meta"`
	}{
		d.Metric,
		d.Timestamp,
		d.Value,
		d.Tags,
		d.Meta,
	})
}

func (d *DataPoint) Clean() error {
	d.Tags = TagSet(NormalizeTags(d.Tags))
	d.Metric = NormalizeMetricKey(d.Metric)

	switch v := d.Value.(type) {
	case string:
		if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			d.Value = i
		} else if f, err := strconv.ParseFloat(v, 64); err == nil {
			d.Value = f
		}
		// else {
		// 	return fmt.Errorf("Unparseable number %v", v)
		// }
	case uint64:
		if v > math.MaxInt64 {
			d.Value = float64(v)
		}
	case *big.Int:
		if bigMaxInt64.Cmp(v) < 0 {
			if f, err := strconv.ParseFloat(v.String(), 64); err == nil {
				d.Value = f
			}
		}
	}
	return nil
}

func (d *DataPoint) GetFlatMetric(tpl string) (string, error) {
	return FlatMetricKeyAndTags(tpl, d.Metric, d.Tags)
}

type TagSet map[string]string

// Copy creates a new TagSet from t.
func (t TagSet) Copy() TagSet {
	n := make(TagSet)
	for k, v := range t {
		n[k] = v
	}
	return n
}

// Merge adds or overwrites everything from o into t and returns t.
func (t TagSet) Merge(o TagSet) TagSet {
	for k, v := range o {
		t[k] = v
	}
	return t
}
