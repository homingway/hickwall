package newcore

import (
	"bufio"
	"bytes"
	"encoding/json"
	"math"
	"math/big"
	"strconv"
	"time"
)

var bigMaxInt64 = big.NewInt(math.MaxInt64)

type DataPoint struct {
	Metric    Metric            `json:"metric"`
	Timestamp time.Time         `json:"timestamp"`
	Value     interface{}       `json:"value"`
	Tags      TagSet            `json:"tags"`
	Meta      map[string]string `json:"meta,omitempty"`
}

type MultiDataPoint []*DataPoint

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

// func (d *DataPoint) GetFlatMetric(tpl string) (string, error) {
// 	return FlatMetricKeyAndTags(tpl, d.Metric, d.Tags)
// }

// AddTS is the same as Add but lets you specify the timestamp
func AddTS(md *MultiDataPoint, name string, ts time.Time, value interface{}, tags TagSet, datatype string, unit string, desc string) {
	// tags := t.Copy()
	// if datatype != metadata.Unknown {
	// 	metadata.AddMeta(name, nil, "datatype", datatype, false)
	// }
	// if unit != "" {
	// 	metadata.AddMeta(name, nil, "unit", unit, false)
	// }
	// if desc != "" {
	// 	metadata.AddMeta(name, tags, "desc", desc, false)
	// }
	if _, present := tags["host"]; !present {
		tags["host"] = GetHostname()
	} else if tags["host"] == "" {
		delete(tags, "host")
	}

	// conf := config.GetRuntimeConf()
	// if conf.Client.Hostname != "" {
	// 	// hostname should be english
	// 	hostname := collectorlib.NormalizeMetricKey(conf.Client.Hostname)
	// 	if hostname != "" {
	// 		tags["host"] = hostname
	// 	}
	// }

	// tags = AddTags.Copy().Merge(tags)

	// d := collectorlib.DataPoint{
	// 	Metric:    name,
	// 	Timestamp: ts,
	// 	Value:     value,
	// 	Tags:      tags,
	// }
	// log.Debugf("DataPoint: %v", d)
	// *md = append(*md, d)
	*md = append(*md, &DataPoint{
		Metric:    NewMetric(name),
		Timestamp: ts,
		Value:     value,
		Tags:      tags,
	})
}

// Add appends a new data point with given metric name, value, and tags. Tags
// may be nil. If tags is nil or does not contain a host key, it will be
// automatically added. If the value of the host key is the empty string, it
// will be removed (use this to prevent the normal auto-adding of the host tag).
func Add(md *MultiDataPoint, name string, value interface{}, t TagSet, datatype string, unit string, desc string) {
	AddTS(md, name, Now(), value, t, datatype, unit, desc)
}
