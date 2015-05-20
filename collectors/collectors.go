package collectors

import (
	"github.com/oliveagle/hickwall/config"
	"github.com/oliveagle/hickwall/newcore"
	"time"
)

// AddTS is the same as Add but lets you specify the timestamp
func AddTS(md *newcore.MultiDataPoint, name string, ts time.Time, value interface{}, t newcore.TagSet, datatype string, unit string, desc string) {
	tags := newcore.AddTags.Copy().Merge(t)

	conf := config.GetRuntimeConf()
	if conf != nil {
		tags = tags.Merge(conf.Client.Tags)
	}

	// if datatype != metadata.Unknown {
	//  metadata.AddMeta(name, nil, "datatype", datatype, false)
	// }
	// if unit != "" {
	//  metadata.AddMeta(name, nil, "unit", unit, false)
	// }
	// if desc != "" {
	//  metadata.AddMeta(name, tags, "desc", desc, false)
	// }

	if _, present := tags["host"]; !present {
		tags["host"] = newcore.GetHostname()
	} else if tags["host"] == "" {
		delete(tags, "host")
	}

	if conf != nil && conf.Client.Hostname != "" {
		// hostname should be english
		hostname := newcore.NormalizeMetricKey(conf.Client.Hostname)
		if hostname != "" {
			tags["host"] = hostname
		}
	}

	*md = append(*md, &newcore.DataPoint{
		Metric:    *newcore.NewMetric(name),
		Timestamp: ts,
		Value:     value,
		Tags:      tags,
	})
}

// Add appends a new data point with given metric name, value, and tags. Tags
// may be nil. If tags is nil or does not contain a host key, it will be
// automatically added. If the value of the host key is the empty string, it
// will be removed (use this to prevent the normal auto-adding of the host tag).
func Add(md *newcore.MultiDataPoint, name string, value interface{}, t newcore.TagSet, datatype string, unit string, desc string) {
	AddTS(md, name, newcore.Now(), value, t, datatype, unit, desc)
}
