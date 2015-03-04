package collectors

import (
	// "fmt"
	"fmt"
	// "github.com/kr/pretty"
	"github.com/oliveagle/go-collectors/datapoint"
	"github.com/oliveagle/go-collectors/metadata"
	"github.com/oliveagle/go-collectors/util"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"sync"
	"time"
)

//TODO, detect duplicated metric keys in configuration

type collector_factory_func func(name string, conf interface{}) Collector

var (
	// collector factories
	collector_factories = make(map[string](collector_factory_func))

	builtin_collectors    []Collector
	customized_collectors []Collector

	metric_keys = make(map[string]bool)

	DefaultFreq = time.Second * 1
	timestamp   = time.Now().Unix()
	tlock       sync.Mutex
	AddTags     datapoint.TagSet
)

type Collector interface {
	Run(chan<- *datapoint.DataPoint)
	Name() string
	Init()
}

/*
false,  already exists
true, added successfully
*/
func AddMetricKey(metric_key string) bool {
	_, ok := metric_keys[metric_key]
	if ok == true {
		return false
	}
	metric_keys[metric_key] = true
	return true
}

// states.key_map[query], _ = ValidateKey(metric_key)
func ValidateKey(key string) (string, error) {
	if AddMetricKey(key) != true {
		err := fmt.Errorf("duplicated metric key: %s", key)
		fmt.Println(err)
		os.Exit(1)
		return "", err
	}
	return key, nil
}

func GetCollectorFactoryByName(name string) (collector_factory_func, bool) {
	factory, ok := collector_factories[name]
	return factory, ok
}

func GetBuiltinCollectorByName(name string) Collector {
	for _, c := range builtin_collectors {
		if c.Name() == name {
			return c
		}
	}
	return nil
}

func GetBuiltinCollectors() []Collector {
	return builtin_collectors
}

func GetCustomizedCollectors() []Collector {
	return customized_collectors
}

func AddCustomizedCollectorByName(factory_name, collector_name string, config interface{}) bool {
	var collector Collector
	factory, ok := GetCollectorFactoryByName(factory_name)
	if ok == true {
		collector = factory(collector_name, config)
		customized_collectors = append(customized_collectors, collector)
	}
	return ok
}

func RemoveAllCustomizedCollectors() {
	customized_collectors = nil
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
	F        func(states interface{}) (datapoint.MultiDataPoint, error)
	Interval time.Duration // default to DefaultFreq
	Enable   func() bool

	name string
	init func()

	states interface{}

	// internal use
	sync.Mutex
	enabled bool
}

func (c *IntervalCollector) Init() {
	if c.init != nil {
		c.init()
	}
}

func (c *IntervalCollector) SetInterval(d time.Duration) {
	c.Interval = d
}

func (c *IntervalCollector) Run(dpchan chan<- *datapoint.DataPoint) {
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
		// pretty.Println("c.Init: ", c)
		// fmt.Println("c.Run ", c)

		next := time.After(interval)
		if c.Enabled() {
			// fmt.Println("Enabled")

			md, err := c.F(c.states)
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

func (c *IntervalCollector) Enabled() bool {
	if c.Enable == nil {
		return true
	}
	c.Lock()
	defer c.Unlock()
	return c.enabled
}

func (c *IntervalCollector) Name() string {
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
