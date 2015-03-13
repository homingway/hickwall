package backends

import (
	"fmt"
	"github.com/influxdb/influxdb/client"
	"github.com/oliveagle/hickwall/config"
	// client088 "github.com/influxdb/influxdb_088/client"
	"github.com/oliveagle/hickwall/utils"
	// "github.com/kr/pretty"
	"github.com/oliveagle/boltq"
	"github.com/oliveagle/go-collectors/datapoint"
	"log"
	// "os"
	"regexp"
	"sync"
	"time"
)

var (
	influxdb_ping_fast_tick = time.Tick(500 * time.Millisecond)
	influxdb_ping_slowtick  = time.Tick(2 * time.Second)
)

var (
	pat_influxdb_version = regexp.MustCompile(`[v]?\d+\.\d+\.\d+[\S]*`)
)

type InfluxdbWriter struct {
	version string
	tick    <-chan time.Time
	tickBkf <-chan time.Time
	mdCh    chan datapoint.MultiDataPoint
	buf     datapoint.MultiDataPoint

	// conf InfluxdbWriterConf
	conf config.Transport_influxdb
	q    *boltq.BoltQ

	iclient InfluxdbClient

	lock_buf      sync.RWMutex
	lock_consume  sync.RWMutex
	lock_backfill sync.RWMutex

	is_consuming   bool
	is_backfilling bool
	is_host_alive  bool

	ping_time_avg int64
}

type InfluxdbWriterConf struct {
	Version        string
	Enabled        bool
	Interval_ms    int
	Max_batch_size int

	// Client Config
	Host     string // for v0.8.8
	URL      string // for v0.9.0
	Username string
	Password string
	Database string

	// Write Config
	RetentionPolicy string
	FlatTemplate    string

	Backfill_enabled              bool
	Backfill_interval_s           int
	Backfill_handsoff             bool
	Backfill_latency_threshold_ms int
	Backfill_cool_down_s          int

	Merge_Requests bool // try best to merge small group of points to no more than max_batch_size
}

func NewInfluxdbWriter(conf config.Transport_influxdb) *InfluxdbWriter {

	version := influxdbParseVersionFromString(conf.Version)

	iclient, err := NewInfluxdbClient(map[string]interface{}{
		"Host":         conf.Host,
		"URL":          conf.URL,
		"Username":     conf.Username,
		"Password":     conf.Password,
		"UserAgent":    "",
		"Database":     conf.Database,
		"FlatTemplate": conf.FlatTemplate,
	}, version)
	if err != nil {
		log.Panicln("cannot create influxdb client: %v", err)
		// log.Printf("cannot create influxdb client: %v", err)
		// os.Exit(1)
	}
	// fmt.Println("InfluxdbClient:", iclient, err)

	//TODO: boltq name should be configurable or automatic generated based on writer's name
	q, err := boltq.NewBoltQ("backend_influxdb.queue", MAX_QUEUE_SIZE, boltq.POP_ON_FULL)
	if err != nil {
		log.Panicf("cannot open backend_influxdb.queue: %v", err)
		// log.Printf("cannot open backend_influxdb.queue: %v", err)
		// os.Exit(1)
	}

	return &InfluxdbWriter{
		version: version,
		conf:    conf,
		tick:    time.Tick(time.Millisecond * time.Duration(conf.Interval_ms)),
		tickBkf: time.Tick(time.Second * time.Duration(conf.Backfill_interval_s)),

		// mdCh must a buffered channel. and if buffer is full. should not write
		// otherwise, program will block. other Tick  within the same goruntime
		// will also be blocked.
		// mdCh: make(chan datapoint.MultiDataPoint, conf.Max_batch_size),
		// we are using `go w.addMD2Buf`, blocking is ok. but it's still better have a buf
		mdCh:    make(chan datapoint.MultiDataPoint, conf.Max_batch_size),
		buf:     datapoint.MultiDataPoint{},
		q:       q,
		iclient: iclient,
	}
}

func (w *InfluxdbWriter) Enabled() bool {
	return w.conf.Enabled
}

func (w *InfluxdbWriter) Close() {
	w.flushToQueue()
}

func (w *InfluxdbWriter) Write(md datapoint.MultiDataPoint) {
	if w.Enabled() == true {
		w.mdCh <- md
	}
}

func (w *InfluxdbWriter) Ping() {
	go func() {
		err_cnt := 0
		fasttick := influxdb_ping_fast_tick
		slowtick := influxdb_ping_slowtick
		tick := fasttick
		failing_pint_cnt := 5
		ping_avg_cnt := 5
		sma := utils.SMA{N: ping_avg_cnt}
		for {
			select {
			case <-tick:
				t, v, err := w.iclient.Ping()
				if err != nil {
					fmt.Println("iclient.Ping: ", err)
					err_cnt += 1
					w.ping_time_avg = 999999999
				} else {

					if v != "" {
						ver := influxdbParseVersionFromString(v)
						if w.iclient.IsCompatibleVersion(ver) == false {
							log.Printf("Error: remote host is incompatible: %s, backend expected version: %s", ver, w.version)
							err_cnt += failing_pint_cnt
							w.ping_time_avg = 999999999
							w.is_host_alive = false
							tick = slowtick
							continue
						}
					}

					if w.is_host_alive == true {
						avg, err := sma.Calc(float64(t.Nanoseconds() / 1000000))
						if err == nil {
							w.ping_time_avg = int64(avg)
						}
					}

					err_cnt = 0
					tick = fasttick
					w.is_host_alive = true
					log.Printf("Fast-PING: resp_time: %v, influxdb ver: %s, q Len: %d, buf size: %d, pingAvg: %d", t, v, w.q.Size(), len(w.buf), w.ping_time_avg)
				}

				if err_cnt > 0 && err_cnt <= failing_pint_cnt {
					log.Printf("Fast-PING: influxdb host is failling %d, q Len: %d, buf size: %d, pingAvg: %d", err_cnt, w.q.Size(), len(w.buf), w.ping_time_avg)
				}

				if err_cnt > failing_pint_cnt {
					w.is_host_alive = false
					tick = slowtick
					log.Printf("SLOW-PING: Wait for influxdb host back online again! q Len: %d, buf size: %d, pingAvg: %d", w.q.Size(), len(w.buf), w.ping_time_avg)
				}

			}
		}
	}()
}

func (w *InfluxdbWriter) Run() {
	go w.Ping()

	for {
		select {
		case md := <-w.mdCh:
			go w.addMD2Buf(md)
		case <-w.tick:
			w.flushToQueue()
			go w.consume()
		case <-w.tickBkf:
			go w.backfill()
		}
	}
}

// --------------------------------------------------------------------
//  Internal Helpers
// --------------------------------------------------------------------

// internal buf, which holds metrics for a very short period. such as 1 second.
func (w *InfluxdbWriter) addMD2Buf(md datapoint.MultiDataPoint) {
	if w.Enabled() == false {
		return
	}

	w.lock_buf.Lock()
	defer w.lock_buf.Unlock()

	for _, p := range md {
		w.buf = append(w.buf, p)
		// log.Println("len(w.buf): ", len(w.buf))
		if len(w.buf) >= w.conf.Max_batch_size {
			// log.Println("make it a batch")
			md1 := datapoint.MultiDataPoint(w.buf[:len(w.buf)])
			// log.Println("addMD2Buf: md: ------------ len: ", len(md1))
			MdPush(w.q, md1)

			w.buf = nil
			// log.Println("w.buf = nil")
		}
	}
}

func (w *InfluxdbWriter) flushToQueue() {
	if w.Enabled() == false {
		return
	}

	w.lock_buf.Lock()
	defer w.lock_buf.Unlock()

	if len(w.buf) > 0 {
		md := datapoint.MultiDataPoint(w.buf[:len(w.buf)])
		// log.Println("flushToQueue: md: ------------ len: ", len(md))
		MdPush(w.q, md)
		w.buf = nil
	}
}

func (w *InfluxdbWriter) consume() {
	if w.Enabled() == false || w.is_consuming == true || w.is_host_alive == false {
		return
	}

	w.lock_consume.Lock()
	defer w.lock_consume.Unlock()
	w.is_consuming = true
	defer func() {
		w.is_consuming = false
	}()

	if w.q.Size() == 0 {
		return
	}

	var md datapoint.MultiDataPoint
	var err error
	if w.conf.Merge_Requests == true {
		md, err = MdPopMany(w.q, w.conf.Max_batch_size)
	} else {
		md, err = MdPop(w.q)
	}

	if err != nil {
		log.Println(err)
		return
	}

	// Do something
	log.Printf(" * md len:%d [influxdb] consuming: boltQ len: %d , mdCh len: %d, buf size: %d\n", len(md), w.q.Size(), len(w.mdCh), len(w.buf))

	err = w.writeMd(md)
	w.check1()

	// when error happened during consume, md will be push back to queue again
	if err != nil {
		log.Printf(" !!! md len:%d -consume- failed, pushback md", len(md))
		MdPush(w.q, md)
	}

}

func (w *InfluxdbWriter) backfill() {
	if w.Enabled() == false || w.is_backfilling == true || w.is_host_alive == false {
		return
	}

	w.lock_backfill.Lock()
	defer w.lock_backfill.Unlock()
	w.is_backfilling = true
	defer func() {
		w.is_backfilling = false
	}()

	// cool down if ping_time_avg break the threshold
	// if w.ping_time_avg >
	if w.conf.Backfill_handsoff == true && w.conf.Backfill_latency_threshold_ms >= 1 && w.conf.Backfill_cool_down_s > 0 && w.ping_time_avg > int64(w.conf.Backfill_latency_threshold_ms) {
		log.Printf(" - backfill is cooling down for %d seconds ----- I'M HOT -----", w.conf.Backfill_cool_down_s)
		time.Sleep(time.Second * time.Duration(w.conf.Backfill_cool_down_s))
		log.Printf(" - backfill is cooling down for %d seconds ----- I'M COOL -----", w.conf.Backfill_cool_down_s)
	}

	if w.conf.Backfill_enabled == true && w.q.Size() > 0 {
		// backfill when boltq is not empty

		var md datapoint.MultiDataPoint
		var err error
		if w.conf.Merge_Requests == true {
			md, err = MdPopManyBottom(w.q, w.conf.Max_batch_size)
		} else {
			md, err = MdPopBottom(w.q)
		}

		if err != nil {
			log.Println(err)
			return
		}

		// do something with backfilling
		log.Printf(" - md len:%d [influxdb] backfilling:, boltQ len: %d\n", len(md), w.q.Size())
		// push back to queue.

		err = w.writeMd(md)
		w.check1()

		if err != nil {
			// push back to queue.
			log.Printf(" !!! md len:%d -backfill- failed, pushback md", len(md))
			MdPush(w.q, md)
		}
	}
}

func (w *InfluxdbWriter) writeMd(md datapoint.MultiDataPoint) error {
	points := []client.Point{}
	for _, p := range md {
		// t, err := client.EpochToTime(p.Timestamp, "n")
		// if err != nil {
		// 	log.Panicln(err)
		// }
		// log.Println(t)
		points = append(points, client.Point{
			Name:      p.Metric,
			Timestamp: p.Timestamp,
			Fields: map[string]interface{}{
				"Value": p.Value,
			},
			Tags: p.Tags, //TODO: Tags
		})
	}
	// log.Println(points)
	write := client.BatchPoints{
		Database:        w.conf.Database,
		RetentionPolicy: w.conf.RetentionPolicy,
		Points:          points,
	}
	res, err := w.iclient.Write(write)
	if err != nil {
		log.Println(" -E- writeMD failed: ", err)
		return err
	}
	if res != nil && res.Err != nil {
		log.Println(" -E- writeMD failed: res.Err: ", res.Err)
		return fmt.Errorf("res.Err: %s", res.Err)
	}
	return err
}

func (w *InfluxdbWriter) check1() {
	// fmt.Println("----- check1 -------")
	db := "metrics"
	q := "select count(Value) from metric1.1"
	res, err := w.iclient.Query(client.Query{
		Command:  q,
		Database: db,
	})
	if err != nil {
		// fmt.Println("err", err)
		return
	}
	if res != nil {
		// fmt.Println("len res.Results: ", len(res.Results))
		for _, r := range res.Results {
			// fmt.Println(r)
			for _, s := range r.Series {
				for _, v := range s.Values {
					// pretty.Println(v)
					if len(v) == 2 {
						log.Printf("-check- queried metric1.1 value Count: %v\n", v[1])
					}
				}
			}
		}
	} else {
		fmt.Println("res == nil ")
	}
}
