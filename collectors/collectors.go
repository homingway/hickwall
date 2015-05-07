package collectors

import (
	// "fmt"
	"fmt"
	log "github.com/oliveagle/seelog"
	// "github.com/kr/pretty"
	"github.com/oliveagle/hickwall/collectorlib"
	"github.com/oliveagle/hickwall/collectorlib/metadata"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"sync"
	"time"

	"github.com/oliveagle/hickwall/config"
)

// func init() {

// 	md, _ := c_hickwall(nil)
// }

//TODO, detect duplicated metric keys in configuration
// type collector_factory_func func(name string, conf interface{}) Collector
type collector_factory_func func(name string, conf interface{}) <-chan Collector

var (
	// collector factories
	collector_factories = make(map[string](collector_factory_func))

	builtin_collectors    []Collector
	customized_collectors []Collector

	metric_keys = make(map[string]bool)

	DefaultFreq = time.Second * 1
	AddTags     collectorlib.TagSet

	timestamp = time.Now()
	tlock     sync.Mutex
)

type Collector interface {
	Run(chan<- collectorlib.MultiDataPoint)
	Name() string
	Init()
	Enabled() bool
	Stop()
}

func init() {
	go func() {
		for t := range time.Tick(time.Second) {
			tlock.Lock()
			timestamp = t
			tlock.Unlock()
		}
	}()
}

func now() (t time.Time) {
	tlock.Lock()
	t = timestamp
	tlock.Unlock()
	return
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

// AddTS is the same as Add but lets you specify the timestamp
func AddTS(md *collectorlib.MultiDataPoint, name string, ts time.Time, value interface{}, t collectorlib.TagSet, rate metadata.RateType, unit string, desc string) {
	tags := t.Copy()
	if rate != metadata.Unknown {
		metadata.AddMeta(name, nil, "rate", rate, false)
	}
	if unit != "" {
		metadata.AddMeta(name, nil, "unit", unit, false)
	}
	if desc != "" {
		metadata.AddMeta(name, tags, "desc", desc, false)
	}
	if host, present := tags["host"]; !present {
		tags["host"] = collectorlib.Hostname
	} else if host == "" {
		delete(tags, "host")
	}

	conf := config.GetRuntimeConf()
	if conf.Hostname != "" {
		// hostname should be english
		hostname := collectorlib.NormalizeMetricKey(conf.Hostname)
		if hostname != "" {
			tags["host"] = hostname
		}
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
func Add(md *collectorlib.MultiDataPoint, name string, value interface{}, t collectorlib.TagSet, rate metadata.RateType, unit string, desc string) {
	AddTS(md, name, now(), value, t, rate, unit, desc)
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

	chstop        chan (int)
	chstop_enable chan (int)
	isrunning     bool
}

func (c *IntervalCollector) Init() {
	c.chstop = make(chan int, 1)

	if c.init != nil {
		c.init()
	}
}

func (c *IntervalCollector) SetInterval(d time.Duration) {
	c.Interval = d
}

func (c *IntervalCollector) Run(dpchan chan<- collectorlib.MultiDataPoint) {
	c.Lock()
	if c.chstop == nil {
		c.chstop = make(chan int, 1)
		c.chstop_enable = make(chan int, 1)
	}
	c.Unlock()

	if c.isrunning == false {
		c.chstop = make(chan int, 1)

		if c.Enable != nil {
			go func() {
			enable_main_loop:
				for {
					next := time.After(time.Minute * 5)
					c.Lock()
					c.enabled = c.Enable()
					c.Unlock()
					<-next
				enable_wait_loop:
					for {
						// read stop chanel
						select {
						case <-c.chstop_enable:
							// fmt.Println("Stop 3 enable_main_loop")
							log.Infof("Stop 3 enable main loop:  %s", c.name)
							break enable_main_loop
						case <-next:
							break enable_wait_loop
						}
					}
				}
			}()
		}

		c.isrunning = true

	main_loop:
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

			// <-next
		wait_loop:
			for {
				// read stop chanel
				select {
				case <-c.chstop:
					c.chstop_enable <- 1
					log.Infof("Stop 2 main loop: %s", c.name)
					break main_loop
				case <-next:
					break wait_loop
				}
			}
		}

		c.isrunning = false
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

func (c *IntervalCollector) Stop() {
	if c.isrunning == true && c.chstop != nil {
		log.Debugf("stopping collector: %s", c.name)
		c.chstop <- 1
	}
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
	RunBuiltinCollectors(mdCh)
	RunCustomizedCollectors(mdCh)
}

func RunBuiltinCollectors(mdCh chan<- collectorlib.MultiDataPoint) {
	for _, c := range builtin_collectors {
		go c.Run(mdCh)
	}

	go run_heartbeat(mdCh)
}

func run_heartbeat(mdCh chan<- collectorlib.MultiDataPoint) {
	for {
		var md collectorlib.MultiDataPoint

		runtime_conf := config.GetRuntimeConf()
		default_interval := time.Second * time.Duration(1)

		interval, err := collectorlib.ParseInterval(runtime_conf.Heartbeat_interval)
		if err != nil {
			log.Errorf("cannot parse interval of heart_beat: %s - %v", runtime_conf.Heartbeat_interval, err)
			interval = default_interval
		} else {
			if interval < default_interval {
				interval = default_interval
			}
		}

		next := time.After(interval)

		tags := AddTags.Copy().Merge(runtime_conf.Tags)
		Add(&md, "hickwall.client.alive", 1, tags, "", "", "")
		log.Debug("Heartbeat")
		mdCh <- md
		<-next
	}
}

func StopBuiltinCollectors() {
	for _, c := range builtin_collectors {
		// fmt.Println("name: ", c.Name())
		c.Stop()
	}
}

// Customized Collectors ---------------------------------------------------------------

func GetCustomizedCollectors() []Collector {
	return customized_collectors
}

func AddCustomizedCollectorByName(factory_name, collector_name string, config interface{}) bool {
	defer log.Flush()

	factory, ok := GetCollectorFactoryByName(factory_name)
	log.Debugf("factory: %s, ok: %v", factory_name, ok)
	if ok == true {
		for collector := range factory(collector_name, config) {
			log.Debugf("collector created with config: %s", collector.Name())
			customized_collectors = append(customized_collectors, collector)
		}

	}
	return ok
}

func RemoveAllCustomizedCollectors() {
	customized_collectors = nil
}

func RunCustomizedCollectors(mdCh chan<- collectorlib.MultiDataPoint) {
	for _, c := range customized_collectors {
		go c.Run(mdCh)
	}
}

func StopCustomizedCollectors() {
	for _, c := range customized_collectors {
		c.Stop()
	}
}

func CreateCustomizedCollectorsFromRuntimeConf() {
	runtime_conf := config.GetRuntimeConf()
	CreateCustomizedCollectorsFromConf(runtime_conf)
}

func CreateCustomizedCollectorsFromConf(runtime_conf *config.RuntimeConfig) {
	defer log.Flush()

	log.Debug("Creating Customized Collectors")

	for i, conf := range runtime_conf.Collector_win_pdh {
		log.Debugf("creating customized collector: win_pdh:, %s,  %v", fmt.Sprintf("c_win_pdh_%d", i), conf)
		AddCustomizedCollectorByName("win_pdh", fmt.Sprintf("c_win_pdh_%d", i), conf)
	}
	log.Debug("Created win_pdh")

	for i, conf := range runtime_conf.Collector_win_wmi {
		log.Debugf("creating customized collector: win_wmi:, %s,  %v", fmt.Sprintf("c_win_wmi_%d", i), conf)
		AddCustomizedCollectorByName("win_wmi", fmt.Sprintf("c_win_wmi_%d", i), conf)
	}

	log.Debug("Created All Customized Collectors")
}
