package collectors

import (
	// "fmt"
	"fmt"
	log "github.com/oliveagle/seelog"
	// "github.com/kr/pretty"
	"github.com/oliveagle/hickwall/collectorlib"
	"github.com/oliveagle/hickwall/collectorlib/metadata"
	"github.com/oliveagle/hickwall/utils"
	"net/http"
	// "os"
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

// type collector_factory_func func(conf interface{}) <-chan Collector

var (
	// collector factories
	collector_factories = make(map[string](collector_factory_func))

	collectors []Collector

	metric_keys = make(map[string]bool)

	DefaultFreq = time.Second * 1
	AddTags     collectorlib.TagSet

	timestamp = time.Now()
	tlock     sync.Mutex
	md_chan   chan collectorlib.MultiDataPoint
)

type Collector interface {
	Run(chan<- collectorlib.MultiDataPoint)
	Name() string
	Init()
	Enabled() bool
	Stop()
	FactoryName() string
}

func SetDataChan(mdChan chan collectorlib.MultiDataPoint) {
	md_chan = mdChan
}

func GetDataChan() chan collectorlib.MultiDataPoint {
	return md_chan
}

func init() {
	go func() {
		for t := range time.Tick(time.Second) {
			tlock.Lock()
			timestamp = t
			tlock.Unlock()
		}
	}()

	md_chan = make(chan collectorlib.MultiDataPoint)
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
		return "", err
	}
	return key, nil
}

func GetCollectorFactoryByName(name string) (collector_factory_func, bool) {
	factory, ok := collector_factories[name]
	return factory, ok
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

	conf := config.GetRuntimeConf().Client
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

	factory_name string

	closeFunc func()
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
				//TODO: memory leak here
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

func (c *IntervalCollector) FactoryName() string {
	return c.factory_name
}

func (c *IntervalCollector) Stop() {
	if c.isrunning == true && c.chstop != nil {
		log.Debugf("stopping collector: %s", c.name)
		c.chstop <- 1

		if c.closeFunc != nil {
			c.closeFunc()
		}
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

func run_heartbeat(mdCh chan<- collectorlib.MultiDataPoint) {
	for {
		var md collectorlib.MultiDataPoint

		// runtime_conf := config.GetRuntimeConf()
		client_conf := config.GetRuntimeConf().Client
		default_interval := time.Second * time.Duration(1)

		interval, err := collectorlib.ParseInterval(client_conf.Heartbeat_interval)
		if err != nil {
			log.Errorf("cannot parse interval of heart_beat: %s - %v", client_conf.Heartbeat_interval, err)
			interval = default_interval
		} else {
			if interval < default_interval {
				interval = default_interval
			}
		}

		next := time.After(interval)

		tags := AddTags.Copy().Merge(client_conf.Tags)
		Add(&md, "hickwall.client.alive", 1, tags, "", "", "")
		log.Debug("Heartbeat")
		mdCh <- md
		<-next
	}
}

// Collectors ---------------------------------------------------------------

func GetCollectors() []Collector {
	return collectors
}

func AddCollector(factory_name, collector_name string, config interface{}) bool {
	defer log.Flush()

	factory, ok := GetCollectorFactoryByName(factory_name)
	log.Debugf("factory: %s, ok: %v", factory_name, ok)
	if ok == true {
		for collector := range factory(collector_name, config) {
			log.Debugf("collector created with config: %s", collector.Name())
			collectors = append(collectors, collector)
		}

	}
	return ok
}

func RemoveAllCollectors() {
	collectors = nil
}

func RunCollectors() error {
	defer utils.Recover_and_log()

	if md_chan != nil {
		go run_heartbeat(md_chan)

		for _, c := range collectors {
			go c.Run(md_chan)
		}

		return nil
	} else {
		return fmt.Errorf("md_chan is nil")
	}
}

func StopCollectors() {
	for _, c := range collectors {
		c.Stop()
	}
}

func CreateCollectorsFromRuntimeConf() {
	runtime_conf := config.GetRuntimeConf()
	CreateCollectorsFromConf(runtime_conf)
}

func CreateCollectorsFromConf(runtime_conf *config.RuntimeConfig) {
	defer log.Flush()

	log.Debug("Creating Customized Collectors")

	AddCollector("win_pdh", "win_pdh", runtime_conf.Collector_win_pdh)
	log.Debug("Created win_pdh")

	AddCollector("win_wmi", "win_wmi", runtime_conf.Collector_win_wmi)

	AddCollector("win_sys", "win_sys", runtime_conf.Collector_win_sys)

	AddCollector("hickwall_client", "hickwall_client", nil)

	log.Debug("Created All Collectors")
}
