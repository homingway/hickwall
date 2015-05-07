package backends

import (
	"fmt"
	"github.com/oliveagle/hickwall/collectorlib"
	"github.com/oliveagle/hickwall/config"
	log "github.com/oliveagle/seelog"
	"strings"
)

var (
	backends []TSWriter

	MAX_BATCH_SIZE = 200
	HttpTimeoutMS  = 500
	MAX_QUEUE_SIZE = int64(100)
)

type TSWriter interface {
	Write(collectorlib.MultiDataPoint)
	Run()
	Close()
	Enabled() bool
	Name() string
}

// func init() {
// }

func GetBackendList() (res []string) {
	for _, bk := range backends {
		res = append(res, bk.Name())
	}
	return
}

func GetBackendByName(name string) (w TSWriter, b bool) {
	for _, bk := range backends {
		if strings.ToLower(bk.Name()) == strings.ToLower(name) {
			return bk, true
		}
	}
	return nil, false
}

func GetBackendByNameVersion(name, version string) (TSWriter, bool) {
	key := strings.Join([]string{name, version}, "-")
	return GetBackendByName(strings.ToLower(key))
}

func WriteToBackends(md collectorlib.MultiDataPoint) {
	for _, bk := range backends {
		if bk.Enabled() == true {
			bk.Write(md)
		}
	}
}

func CloseBackends() {
	for _, bk := range backends {
		bk.Close()
		log.Debug("Closed Backend ", bk.Name())
	}
}

func RunBackends() {
	for _, bk := range backends {
		if bk.Enabled() == true {
			log.Debug("Backend is Running: ", bk.Name())
			go bk.Run()
		} else {
			log.Debug("Backend is Not Running: ", bk.Name())
		}
	}
}

func RemoveAllBackends() {
	backends = nil
}

func CreateBackendsFromRuntimeConf() {
	log.Debug("Create backends from runtime config")
	runtime_conf := config.GetRuntimeConf()
	CreateBackendsFromConf(runtime_conf)
}

func CreateBackendsFromConf(runtime_conf *config.RuntimeConfig) {
	defer log.Flush()

	log.Debug("creating backends from conf")

	backends = append(backends, NewStdoutWriter("stdout", runtime_conf.Transport_stdout))
	log.Debug("stdout backend created")

	backends = append(backends, NewFileWriter("file", runtime_conf.Transport_file))
	log.Debug("file backend created")

	// influxdb backends
	for _, iconf := range runtime_conf.Transport_influxdb {
		bkname := fmt.Sprintf("influxdb-%s", influxdbParseVersionFromString(iconf.Version))
		log.Debugf("Creating backend: %s", bkname)
		bk, err := NewInfluxdbWriter(bkname, iconf)
		if err != nil {
			log.Errorf("create backend failed: %v ", err)
			continue
		}
		backends = append(backends, bk)
	}

	log.Debug("all backends created")
}
