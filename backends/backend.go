package backends

import (
	"fmt"
	"github.com/oliveagle/go-collectors/datapoint"
	"github.com/oliveagle/hickwall/config"
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
}

func init() {

	fmt.Println("config: ", config.Conf.Transport_influxdb)

	stdConf := StdoutWriterConf{
		Enabled:           true,
		Max_batch_size:    MAX_BATCH_SIZE,
		Interval:          time.Millisecond * time.Duration(1000),
		Backfill_enabled:  true,
		Backfill_interval: time.Millisecond * time.Duration(200),
	}
	backends["stdout"] = NewStdoutWriter(stdConf)

	// influxdb backends
	for _, iconf := range config.Conf.Transport_influxdb {
		backends[fmt.Sprintf(
			"influxdb-%s",
			influxdbParseVersionFromString(iconf.Version),
		)] = NewInfluxdbWriter(iconf)
	}
}

func GetBackendList() (res []string) {
	for key, _ := range backends {
		res = append(res, key)
	}
	return
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
