package collectors

import (
	// "fmt"
	"fmt"
	// log "github.com/oliveagle/seelog"
	// "github.com/kr/pretty"
	"github.com/oliveagle/hickwall/collectorlib"
	"github.com/oliveagle/hickwall/collectorlib/metadata"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"sync"
	"time"
)

// func init() {

// 	md, _ := c_hickwall(nil)
// }

//TODO, detect duplicated metric keys in configuration
type collector_factory_func func(name string, conf interface{}) Collector

var (
	// collector factories
	collector_factories = make(map[string](collector_factory_func))

	builtin_collectors    []Collector
	customized_collectors []Collector

	metric_keys = make(map[string]bool)

	DefaultFreq = time.Second * 1
	AddTags     collectorlib.TagSet
)

type Collector interface {
	Run(chan<- collectorlib.MultiDataPoint)
	Name() string
	Init()
	Enabled() bool
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

// AddTS is the same as Add but lets you specify the timestamp
func AddTS(md *collectorlib.MultiDataPoint, name string, ts time.Time, value interface{}, t collectorlib.TagSet, rate metadata.RateType, unit metadata.Unit, desc string) {
	tags := t.Copy()
	if rate != metadata.Unknown {
		metadata.AddMeta(name, nil, "rate", rate, false)
	}
	if unit != metadata.None {
		metadata.AddMeta(name, nil, "unit", unit, false)
	}
	// if desc != "" {
	// 	metadata.AddMeta(name, tags, "desc", desc, false)
	// }
	if host, present := tags["host"]; !present {
		tags["host"] = collectorlib.Hostname
	} else if host == "" {
		delete(tags, "host")
	}
	tags = AddTags.Copy().Merge(tags)
	d := collectorlib.DataPoint{
		Metric:    name,
		Timestamp: ts,
		Value:     value,
		Tags:      tags,
	}
	// log.Debugf("DataPoint: %v", d)
	*md = append(*md, d)
}

// Add appends a new data point with given metric name, value, and tags. Tags
// may be nil. If tags is nil or does not contain a host key, it will be
// automatically added. If the value of the host key is the empty string, it
// will be removed (use this to prevent the normal auto-adding of the host tag).
func Add(md *collectorlib.MultiDataPoint, name string, value interface{}, t collectorlib.TagSet, rate metadata.RateType, unit metadata.Unit, desc string) {
	AddTS(md, name, time.Now(), value, t, rate, unit, desc)
}

type IntervalCollector struct {
	F        func(states interface{}) (collectorlib.MultiDataPoint, error)
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

func (c *IntervalCollector) Run(dpchan chan<- collectorlib.MultiDataPoint) {
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

		next := time.After(interval)
		if c.Enabled() {

			md, err := c.F(c.states)

			if err != nil {
				fmt.Errorf("%v: %v", c.Name(), err)
			}
			dpchan <- md
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

func RunAllCollectors(mdCh chan<- collectorlib.MultiDataPoint) {
	for _, c := range builtin_collectors {
		go c.Run(mdCh)
	}
	for _, c := range customized_collectors {
		go c.Run(mdCh)
	}
}
