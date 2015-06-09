package newcore

import (
	"bufio"
	"bytes"
	"encoding/json"
	// "errors"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"sync"
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

var ppFree = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

//TODO: test to make sure buffer don't return heading zero bytes. e.g. [00 00 00 13 42]
// json.Encoding will have too much overhead.
func (d *DataPoint) Json() []byte {
	d.Clean()

	buffer := ppFree.Get().(*bytes.Buffer)
	buffer.WriteString(`{"metric":"`)
	buffer.WriteString(string(d.Metric))
	buffer.WriteString(`","timestamp":"`)
	buffer.WriteString(d.Timestamp.Format(time.RFC3339))
	buffer.WriteString(`","value":`)
	buffer.WriteString(d.value2string())

	//TODO: handle tags
	//TODO: handle meta

	buffer.WriteString(`}`)
	// s := buffer.String()
	s := buffer.Bytes()
	buffer.Reset()
	ppFree.Put(buffer)
	d.length = len(s)
	return s
}

func (d *DataPoint) value2string() string {
	switch t := d.Value.(type) {
	case bool, int, int8, int16, int32, int64, float32, float64:
		return fmt.Sprintf("%v", d.Value)
	default:
		return fmt.Sprintf(`"%v"`, t)
	}
}

func (d *DataPoint) MarshalJSON2String() (string, error) {
	d.Clean()

	// res, err := json.Marshal(struct {
	//  Metric    string            `json:"metric"`
	//  Timestamp time.Time         `json:"timestamp"`
	//  Value     interface{}       `json:"value"`
	//  Tags      TagSet            `json:"tags"`
	//  Meta      map[string]string `json:"meta"`
	// }{
	//  d.Metric,
	//  d.Timestamp,
	//  d.Value,
	//  d.Tags,
	//  d.Meta,
	// })

	// res, err := json.Marshal(d)
	// v := string(res)
	// res = nil

	var b bytes.Buffer
	defer b.Reset()

	writer := bufio.NewWriter(&b)

	enc := json.NewEncoder(writer)

	err := enc.Encode(d)
	writer.Flush()

	v := b.String()

	return v, err
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
	return d.Json(), nil
}

func (d *DataPoint) Length() int {
	if d.length <= 0 {
		d.Encode() // have to Encode.
	}
	return d.length
}
