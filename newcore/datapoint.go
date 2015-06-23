package newcore

import (
	// "bufio"
	// "bytes"
	"encoding/json"
	// "errors"
	"fmt"
	"math"
	"math/big"
	"strconv"
	// "sync"
	"time"
)

var (
	bigMaxInt64 = big.NewInt(math.MaxInt64)
	// timestamp_layout = time.RFC3339
	// "Jan 2, 2006 at 3:04pm (MST)"
)

type DataPoint struct {
	Metric    Metric            `json:"metric"`
	Timestamp time.Time         `json:"timestamp"`
	Value     interface{}       `json:"value"`
	Tags      TagSet            `json:"tags"`
	Meta      map[string]string `json:"meta,omitempty"`

	length int
}

func NewDPFromJson(content []byte) (dp DataPoint, err error) {
	err = json.Unmarshal(content, &dp)
	return
}

func NewDP(prefix, metric string, value interface{}, tags TagSet, datatype string, unit string, desc string) DataPoint {
	return NewDataPoint(
		fmt.Sprintf("%s.%s", prefix, metric),
		value,
		Now(),
		tags,
		datatype,
		unit,
		desc,
	)
}

//TODO: unittest NewDataPoint
func NewDataPoint(metric string, value interface{}, ts time.Time, t TagSet, datatype string, unit string, desc string) DataPoint {
	tags := AddTags.Copy().Merge(t)

	if _, present := tags["host"]; !present {
		tags["host"] = GetHostname()
	} else if tags["host"] == "" {
		delete(tags, "host")
	}

	return DataPoint{
		Metric:    Metric(metric),
		Timestamp: ts,
		Value:     value,
		Tags:      tags,
	}
}

type MultiDataPoint []DataPoint

// var ppFree = sync.Pool{
// 	New: func() interface{} {
// 		return new(bytes.Buffer)
// 	},
// }

func (d *DataPoint) MarshalJSON() ([]byte, error) {

	res, err := json.Marshal(struct {
		Metric    Metric            `json:"metric"`
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
	d.length = len(res)
	return res, err
}

func (d *DataPoint) Clean() error {
	d.Tags = TagSet(NormalizeTags(d.Tags))
	// d.Metric = NormalizeMetricKey(string(d.Metric))
	d.Metric = Metric(d.Metric.Clean())

	switch v := d.Value.(type) {
	case string:
		if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			d.Value = i
		} else if f, err := strconv.ParseFloat(v, 64); err == nil {
			d.Value = f
		}
		// else {
		//  return fmt.Errorf("Unparseable number %v", v)
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

// fulfill sarama kafka Encoder interface{}
func (d *DataPoint) Encode() ([]byte, error) {
	return d.MarshalJSON()
}

// fulfill sarama kafka Encoder interface{}
func (d *DataPoint) Length() int {
	if d.length <= 0 {
		d.Encode() // have to Encode.
	}
	return d.length
}
