package collectorlib

import (
	"fmt"
	"github.com/oliveagle/go-collectors/datapoint"
	"github.com/oliveagle/go-collectors/metadata"
	"github.com/oliveagle/go-collectors/util"
	"net/http"
	"reflect"
	"runtime"
	"sync"
	"time"
)

var (
	collectors []Collector

	DefaultFreq = time.Second * 1
	timestamp   = time.Now().Unix()
	tlock       sync.Mutex
	AddTags     datapoint.TagSet
)

type Collector interface {
	Run(chan<- *datapoint.DataPoint)
	Name() string
	Init()
	SetConfig(config interface{})
	// Close()
}

// Builtin Collectors

// Add Collector from Configuration

func AddCollector(c *Collector) {
	collectors = append(collectors, *c)
}

func GetAllCollectors() []Collector {
	return collectors
}

func init() {
	go func() {
		for t := range time.Tick(time.Second) {
			tlock.Lock()
			timestamp = t.Unix()
			tlock.Unlock()
		}
	}()
}

func now() (t int64) {
	tlock.Lock()
	t = timestamp
	tlock.Unlock()
	return
}

// AddTS is the same as Add but lets you specify the timestamp
func AddTS(md *datapoint.MultiDataPoint, name string, ts int64, value interface{}, t datapoint.TagSet, rate metadata.RateType, unit metadata.Unit, desc string) {
	tags := t.Copy()
	if rate != metadata.Unknown {
		metadata.AddMeta(name, nil, "rate", rate, false)
	}
	if unit != metadata.None {
		metadata.AddMeta(name, nil, "unit", unit, false)
	}
	if desc != "" {
		metadata.AddMeta(name, tags, "desc", desc, false)
	}
	if host, present := tags["host"]; !present {
		tags["host"] = util.Hostname
	} else if host == "" {
		delete(tags, "host")
	}
	tags = AddTags.Copy().Merge(tags)
	d := datapoint.DataPoint{
		Metric:    name,
		Timestamp: ts,
		Value:     value,
		Tags:      tags,
	}
	*md = append(*md, &d)
}

// Add appends a new data point with given metric name, value, and tags. Tags
// may be nil. If tags is nil or does not contain a host key, it will be
// automatically added. If the value of the host key is the empty string, it
// will be removed (use this to prevent the normal auto-adding of the host tag).
func Add(md *datapoint.MultiDataPoint, name string, value interface{}, t datapoint.TagSet, rate metadata.RateType, unit metadata.Unit, desc string) {
	AddTS(md, name, now(), value, t, rate, unit, desc)
}

type IntervalCollector struct {
	F        func() (datapoint.MultiDataPoint, error)
	Interval time.Duration // default to DefaultFreq
	Enable   func() bool
	name     string
	// init     func(c *IntervalCollector, config interface{})
	init func()

	config interface{}
	// internal use
	sync.Mutex
	enabled bool
}

func NewIntervalCollector(
	name string,
	init func(),
	F func() (datapoint.MultiDataPoint, error),
	Enable func() bool,
) Collector {
	return IntervalCollector{
		F:      F,
		Enable: Enable,
		name:   name,
		init:   init,
	}
}

func (c IntervalCollector) SetConfig(config interface{}) {
}

func (c IntervalCollector) Init() {
	if c.init != nil {
		c.init()
	}
}

func (c IntervalCollector) Run(dpchan chan<- *datapoint.DataPoint) {
	if c.Enable != nil {
		go func() {
			for {
				next := time.After(time.Minute * 5)
				c.Lock()
				c.enabled = c.Enable()
				c.Unlock()
				<-next
			}
		}()
	}
	for {

		interval := c.Interval
		if interval == 0 {
			interval = DefaultFreq
		}
		// fmt.Println(time.Now(), "for interval: ", interval)

		next := time.After(interval)
		if c.Enabled() {
			// fmt.Println("Enabled")

			md, err := c.F()
			// fmt.Println(md, err)

			if err != nil {
				// slog.Errorf("%v: %v", c.Name(), err)
				fmt.Errorf("%v: %v", c.Name(), err)
			}
			for _, dp := range md {
				dpchan <- dp
			}
		}
		<-next
	}
}

func (c IntervalCollector) Enabled() bool {
	if c.Enable == nil {
		return true
	}
	c.Lock()
	defer c.Unlock()
	return c.enabled
}

func (c IntervalCollector) Name() string {
	if c.name != "" {
		return c.name
	}
	v := runtime.FuncForPC(reflect.ValueOf(c.F).Pointer())
	return v.Name()
}

func enableURL(url string) func() bool {
	return func() bool {
		resp, err := http.Get(url)
		if err != nil {
			return false
		}
		resp.Body.Close()
		return resp.StatusCode == 200
	}
}
