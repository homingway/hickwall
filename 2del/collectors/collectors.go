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

	"container/list"
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

	timestamp        = time.Now()
	tlock            sync.Mutex
	md_chan          chan collectorlib.MultiDataPoint
	chstop_heartbeat chan bool

	// running_collecotrs [](chan<- bool)
	running_collecotrs *list.List
)

// type Collector interface {
// 	Run(chan<- collectorlib.MultiDataPoint) chan<- bool
// 	Name() string
// 	Init()
// 	IsEnabled() bool
// 	FactoryName() string
// }

// func SetDataChan(mdChan chan collectorlib.MultiDataPoint) {
// 	md_chan = mdChan
// }

// func GetDataChan() chan collectorlib.MultiDataPoint {
// 	return md_chan
// }

// func init() {
// 	go func() {
// 		for t := range time.Tick(time.Second) {
// 			tlock.Lock()
// 			timestamp = t
// 			tlock.Unlock()
// 		}
// 	}()

// 	// md_chan = make(chan collectorlib.MultiDataPoint)
// 	md_chan = make(chan collectorlib.MultiDataPoint, 1000)

// 	chstop_heartbeat = make(chan bool)

// 	running_collecotrs = list.New()
// }

// func now() (t time.Time) {
// 	tlock.Lock()
// 	t = timestamp
// 	tlock.Unlock()
// 	return
// }

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

// // AddTS is the same as Add but lets you specify the timestamp
// func AddTS(md *collectorlib.MultiDataPoint, name string, ts time.Time, value interface{}, tags collectorlib.TagSet, rate metadata.RateType, unit string, desc string) {
// 	// tags := t.Copy()
// 	// if rate != metadata.Unknown {
// 	// 	metadata.AddMeta(name, nil, "rate", rate, false)
// 	// }
// 	// if unit != "" {
// 	// 	metadata.AddMeta(name, nil, "unit", unit, false)
// 	// }
// 	// if desc != "" {
// 	// 	metadata.AddMeta(name, tags, "desc", desc, false)
// 	// }
// 	if _, present := tags["host"]; !present {
// 		tags["host"] = collectorlib.Hostname
// 	} else if tags["host"] == "" {
// 		delete(tags, "host")
// 	}

// 	conf := config.GetRuntimeConf()
// 	if conf.Client.Hostname != "" {
// 		// hostname should be english
// 		hostname := collectorlib.NormalizeMetricKey(conf.Client.Hostname)
// 		if hostname != "" {
// 			tags["host"] = hostname
// 		}
// 	}
// 	// tags = AddTags.Copy().Merge(tags)

// 	// d := collectorlib.DataPoint{
// 	// 	Metric:    name,
// 	// 	Timestamp: ts,
// 	// 	Value:     value,
// 	// 	Tags:      tags,
// 	// }
// 	// log.Debugf("DataPoint: %v", d)
// 	// *md = append(*md, d)
// 	*md = append(*md, &collectorlib.DataPoint{
// 		Metric:    name,
// 		Timestamp: ts,
// 		Value:     value,
// 		Tags:      tags,
// 	})
// }

// // Add appends a new data point with given metric name, value, and tags. Tags
// // may be nil. If tags is nil or does not contain a host key, it will be
// // automatically added. If the value of the host key is the empty string, it
// // will be removed (use this to prevent the normal auto-adding of the host tag).
// func Add(md *collectorlib.MultiDataPoint, name string, value interface{}, t collectorlib.TagSet, rate metadata.RateType, unit string, desc string) {
// 	AddTS(md, name, now(), value, t, rate, unit, desc)
// }

// type IntervalCollector struct {
// 	F          func(states interface{}) (collectorlib.MultiDataPoint, error)
// 	Interval   time.Duration // default to DefaultFreq
// 	EnableFunc func() bool

// 	name string
// 	init func()

// 	states interface{}

// 	// internal use
// 	sync.Mutex
// 	enabled  bool
// 	stopping bool

// 	done chan bool

// 	factory_name string

// 	closeFunc func()
// }

// func (c *IntervalCollector) Init() {
// 	if c.init != nil {
// 		c.init()
// 	}
// }

// func (c *IntervalCollector) SetInterval(d time.Duration) {
// 	c.Interval = d
// }

// func (c *IntervalCollector) loop(dpchan chan<- collectorlib.MultiDataPoint, done chan bool) {
// 	defer utils.Recover_and_log()

// 	// while reloading configuration, consumers maybe closed before
// 	// internal_buf := make(chan collectorlib.MultiDataPoint, 1000)

// 	log.Infof("START IntervalCollector! %08x, name: %s,  wait on done: %08x", &c, c.Name(), &done)

// 	interval := c.Interval
// 	if interval == 0 {
// 		interval = DefaultFreq
// 	}

// 	tick := time.Tick(interval)
// 	tick_enabler := time.Tick(time.Minute * 5)

// collector_loop:
// 	for {
// 		select {
// 		case <-tick:
// 			if c.IsEnabled() {
// 				//TODO: memory leak here
// 				// if c.stopping == true {
// 				// 	break
// 				// }

// 				md, err := c.F(c.states)
// 				if err != nil {
// 					fmt.Errorf("%v: %v", c.Name(), err)
// 					break
// 				}

// 				dpchan <- md

// 				// select {
// 				// case dpchan <- md:
// 				// 	continue
// 				// case <-done:
// 				// 	log.Infof("break collector_loop: %08x - in <-tick", &c)
// 				// 	// this did happened occasionally.
// 				// 	// win > try_reload_config (master) 17:04:58 $ ./clean.sh && go run try_reload_config.go | grep "<-tick"
// 				// 	// 2015-05-13T17:05:13.97 CST [Info] collectors.go:240(loop) break collector_loop: c08204c248 - in <-tick
// 				// 	// 2015-05-13T17:05:17.96 CST [Info] collectors.go:240(loop) break collector_loop: c08204c1b0 - in <-tick
// 				// 	// 2015-05-13T17:05:40.98 CST [Info] collectors.go:240(loop) break collector_loop: c08204c180 - in <-tick
// 				// 	// 2015-05-13T17:05:43.98 CST [Info] collectors.go:240(loop) break collector_loop: c08204c1c8 - in <-tick
// 				// 	break collector_loop
// 				// }
// 			}
// 		case <-tick_enabler:
// 			if c.EnableFunc != nil {
// 				c.Lock()
// 				c.enabled = c.EnableFunc()
// 				c.Unlock()
// 			}
// 		case <-done:
// 			// log.Infof("break collector_loop: %08x - in <-done", &c)
// 			break collector_loop
// 		}
// 	}
// 	log.Infof("STOP IntervalCollector! %08x", &c)
// }

// func (c *IntervalCollector) sending(dpchan chan<- collectorlib.MultiDataPoint, md collectorlib.MultiDataPoint) {
// 	// c.Lock()
// 	// defer c.Unlock()
// 	if c.stopping == false {
// 		dpchan <- md
// 	}
// }

// func (c *IntervalCollector) setStopping(b bool) {
// 	c.Lock()
// 	defer c.Unlock()
// 	c.stopping = b
// }

// func (c *IntervalCollector) Run(dpchan chan<- collectorlib.MultiDataPoint) chan<- bool {
// 	c.Lock()
// 	defer c.Unlock()

// 	done := make(chan bool)
// 	go c.loop(dpchan, done)
// 	return done
// }

// func (c *IntervalCollector) IsEnabled() bool {
// 	if c.EnableFunc == nil {
// 		return true
// 	}
// 	c.Lock()
// 	defer c.Unlock()
// 	return c.enabled
// }

// func (c *IntervalCollector) Name() string {
// 	if c.name != "" {
// 		return c.name
// 	}
// 	v := runtime.FuncForPC(reflect.ValueOf(c.F).Pointer())
// 	return v.Name()
// }

// func (c *IntervalCollector) FactoryName() string {
// 	return c.factory_name
// }

// func enableURL(url string) func() bool {
// 	return func() bool {
// 		resp, err := http.Get(url)
// 		if err != nil {
// 			return false
// 		}
// 		resp.Body.Close()
// 		return resp.StatusCode == 200
// 	}
// }

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
	// running_collecotrs = nil
}

func RunCollectors() error {
	defer utils.Recover_and_log()

	if md_chan != nil {

		// for _, c := range collectors {
		// 	running_collecotrs = append(running_collecotrs, c.Run(md_chan))
		// }
		for _, c := range collectors {
			running_collecotrs.PushBack(c.Run(md_chan))
		}

		return nil
	} else {
		return fmt.Errorf("md_chan is nil")
	}
}

func StopCollectors() {

	log.Info("StopCollectors START")
	for el := running_collecotrs.Front(); el != nil; el = running_collecotrs.Front() {
		done := el.Value.(chan<- bool)
		log.Infof("sending true to done: %08x", &done)
		done <- true
		// done <- true
		running_collecotrs.Remove(el)
	}
	log.Info("StopCollectors DONE")
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

//----------------------------------- heart beat ----------------------------------------------------------

// func run_heartbeat(mdCh chan<- collectorlib.MultiDataPoint, done <-chan bool) {
// 	var (
// 		md               collectorlib.MultiDataPoint
// 		default_interval = time.Second * time.Duration(1)
// 	)
// 	client_conf := config.GetRuntimeConf().Client

// 	interval, err := collectorlib.ParseInterval(client_conf.Heartbeat_interval)
// 	if err != nil {
// 		log.Errorf("cannot parse interval of heart_beat: %s - %v", client_conf.Heartbeat_interval, err)
// 		interval = default_interval
// 	} else {
// 		if interval < default_interval {
// 			interval = default_interval
// 		}
// 	}

// 	tick := time.Tick(interval)

// loop:
// 	for {
// 		select {
// 		case <-tick:
// 			tags := AddTags.Copy().Merge(client_conf.Tags)
// 			Add(&md, "hickwall.client.alive", 1, tags, "", "", "")
// 			log.Debug("Heartbeat")
// 			mdCh <- md
// 			md = nil
// 		case <-done:
// 			log.Info("heartbeat stopped")
// 			break loop
// 		}
// 	}
// }

// var hb *HeartBeater

// func StartHeartBeat() {
// 	if hb == nil {
// 		hb = NewHeartBeater()
// 	}
// 	hb.Start()
// }

// func StopHeartBeat() {
// 	if hb == nil {
// 		hb = NewHeartBeater()
// 	}
// 	hb.Stop()
// }

// type HeartBeater struct {
// 	done     chan bool
// 	running  bool
// 	interval time.Duration
// }

// func NewHeartBeater() *HeartBeater {
// 	var (
// 		default_interval = time.Second * time.Duration(1)
// 		interval         = default_interval
// 	)

// 	conf := config.GetRuntimeConf()
// 	if conf != nil {
// 		client_conf := conf.Client

// 		interval, err := collectorlib.ParseInterval(client_conf.Heartbeat_interval)
// 		if err != nil {
// 			log.Errorf("cannot parse interval of heart_beat: %s - %v", client_conf.Heartbeat_interval, err)
// 			interval = default_interval
// 		} else {
// 			if interval < default_interval {
// 				interval = default_interval
// 			}
// 		}
// 	} else {
// 		interval = default_interval
// 	}

// 	return &HeartBeater{
// 		done:     make(chan bool),
// 		interval: interval,
// 	}
// }

// func (h *HeartBeater) IsRunning() bool {
// 	return h.running
// }

// func (h *HeartBeater) Start() {
// 	if h.IsRunning() == false {

// 		go func() {
// 			var md collectorlib.MultiDataPoint

// 			client_conf := config.GetRuntimeConf().Client

// 			log.Info("*HeartBeater START")
// 			tick := time.Tick(h.interval)

// 			mdCh := GetDataChan()

// 		loop:
// 			for {
// 				select {
// 				case <-tick:
// 					tags := AddTags.Copy().Merge(client_conf.Tags)
// 					Add(&md, "hickwall.client.alive", 1, tags, "", "", "")
// 					log.Debug("Heartbeat")
// 					mdCh <- md
// 					md = nil
// 				case <-h.done:
// 					log.Info("*HeartBeater STOP 2")
// 					break loop
// 				}
// 			}

// 			log.Info("HeartBeater: gorotuine finished.")
// 		}()
// 		h.running = true
// 	}
// }

// func (h *HeartBeater) Stop() {
// 	if h.IsRunning() == true {
// 		log.Info("*HeartBeater STOP 1")
// 		h.done <- true
// 		h.running = false
// 	}
// }
