package backends

import (
	"fmt"
	"github.com/oliveagle/go-collectors/datapoint"
	"strings"
	"time"
)

var (
	backends       = map[string]TSWriter{}
	MAX_BATCH_SIZE = 200
	HttpTimeoutMS  = 500
	MAX_QUEUE_SIZE = int64(100)
)

type TSWriter interface {
	Write(datapoint.MultiDataPoint)
	Run()
	Close()
	Enabled() bool

	// GetTsWriterConf() *TSWriterConf
}

// type TSWriterConf struct {
// 	Interval                      time.Duration
// 	Backfill_interval             time.Duration
// 	Max_batch_size                int
// 	Backfill_enabled              bool
// 	Backfill_handsoff             bool
// 	Backfill_latency_threshold_ms int
// 	Backfill_cool_down_second     int
// }

func init() {
	// tsconf := TSWriterConf{

	// // Backfill_enabled:  false,

	// }
	stdConf := StdoutWriterConf{
		Enabled:           true,
		Max_batch_size:    MAX_BATCH_SIZE,
		Interval:          time.Millisecond * time.Duration(1000),
		Backfill_enabled:  true,
		Backfill_interval: time.Millisecond * time.Duration(200),
	}
	backends["stdout"] = NewStdoutWriter(stdConf)

	influxdbConf_v090 := InfluxdbWriterConf{
		Version:        "v0.9.0-rc7",
		Enabled:        true,
		Max_batch_size: MAX_BATCH_SIZE,
		Interval_ms:    1000,

		URL:             "http://192.168.59.103:8086/write",
		Username:        "root",
		Password:        "root",
		Database:        "metrics",
		RetentionPolicy: "p1",

		Backfill_enabled:    true,
		Backfill_interval_s: 1,

		Backfill_handsoff:             true,
		Backfill_latency_threshold_ms: 10,
		Backfill_cool_down_s:          5,
		Merge_Requests:                false,
	}
	backends["influxdb-0.9.0-rc7"] = NewInfluxdbWriter(influxdbConf_v090)

	influxdbConf_v088 := InfluxdbWriterConf{
		Version:        "v0.8.8",
		Enabled:        true,
		Max_batch_size: MAX_BATCH_SIZE,
		Interval_ms:    1000,

		Host:     "192.168.59.103:8086",
		URL:      "http://192.168.59.103:8086/write",
		Username: "root",
		Password: "root",
		Database: "metrics",

		Backfill_enabled:    true,
		Backfill_interval_s: 1,

		Backfill_handsoff:             true,
		Backfill_latency_threshold_ms: 10,
		Backfill_cool_down_s:          5,
		Merge_Requests:                false,
	}
	backends["influxdb-0.8.8"] = NewInfluxdbWriter(influxdbConf_v088)
}

func GetBackendByName(name string) (w TSWriter, b bool) {
	w, b = backends[strings.ToLower(name)]
	return
}

func GetBackendByNameVersion(name, version string) (w TSWriter, b bool) {
	key := strings.Join([]string{name, version}, "-")
	w, b = backends[strings.ToLower(key)]
	return
}

func WriteToBackends(md datapoint.MultiDataPoint) {
	for _, backend := range backends {
		if backend.Enabled() == true {
			backend.Write(md)
		}
	}
}

func CloseBackends() {
	for _, backend := range backends {
		backend.Close()
	}
}

func RunBackends() {
	for key, backend := range backends {
		fmt.Printf("backend: %s ", key)
		if backend.Enabled() == true {
			fmt.Println("Running")
			go backend.Run()
		} else {
			fmt.Println("Not Running")
		}
	}
}
