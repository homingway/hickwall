package backends

import (
	"fmt"
	"github.com/oliveagle/hickwall/collectorlib"
	"github.com/oliveagle/hickwall/config"
	log "github.com/oliveagle/seelog"
	"strings"
)

var (
	backends       = map[string]TSWriter{}
	MAX_BATCH_SIZE = 200
	HttpTimeoutMS  = 500
	MAX_QUEUE_SIZE = int64(100)
)

type TSWriter interface {
	Write(collectorlib.MultiDataPoint)
	Run()
	Close()
	Enabled() bool
}

func init() {
	backends["stdout"] = NewStdoutWriter(config.Conf.Transport_stdout)

	// influxdb backends
	for _, iconf := range config.Conf.Transport_influxdb {
		bkname := fmt.Sprintf("influxdb-%s", influxdbParseVersionFromString(iconf.Version))
		bk, err := NewInfluxdbWriter(iconf)
		if err != nil {
			log.Criticalf("create backend failed: %v ", err)
			continue
		}
		backends[bkname] = bk
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

func WriteToBackends(md collectorlib.MultiDataPoint) {
	for key, backend := range backends {
		if backend.Enabled() == true {
			log.Debug("Backend.Write.Endabled: ", key)
			backend.Write(md)
		}
		// else {
		// 	log.Debug("Backend.Write.Disabled: ", key)
		// }
	}
}

func CloseBackends() {
	for key, backend := range backends {
		backend.Close()
		log.Debug("Closed Backend ", key)
	}
}

func RunBackends() {
	for key, backend := range backends {
		if backend.Enabled() == true {
			log.Debug("Backend is Running: ", key)
			go backend.Run()
		} else {
			log.Debug("Backend is Not Running: ", key)

		}
	}
}
