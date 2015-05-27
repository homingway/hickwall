package backends

import (
	"fmt"
	"github.com/oliveagle/hickwall/collectorlib"
	"github.com/oliveagle/hickwall/config"
	"github.com/oliveagle/hickwall/utils"
	log "github.com/oliveagle/seelog"
	"strings"

	"sync"
)

var (
	backends []TSWriter

	MAX_BATCH_SIZE = 200
	HttpTimeoutMS  = 500
	MAX_QUEUE_SIZE = int64(100)

	mutex sync.Mutex
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
	mutex.Lock()
	defer mutex.Unlock()

	for _, bk := range backends {
		res = append(res, bk.Name())
	}
	return
}

func GetBackendByName(name string) (w TSWriter, b bool) {
	mutex.Lock()
	defer mutex.Unlock()

	for _, bk := range backends {
		if strings.ToLower(bk.Name()) == strings.ToLower(name) {
			return bk, true
		}
	}
	return nil, false
}

func GetBackendByNameVersion(name, version string) (TSWriter, bool) {
	mutex.Lock()
	defer mutex.Unlock()

	key := strings.Join([]string{name, version}, "-")
	return GetBackendByName(strings.ToLower(key))
}

func WriteToBackends(md collectorlib.MultiDataPoint) {
	mutex.Lock()
	defer mutex.Unlock()
	// log.Info("write to backends ")

	for _, bk := range backends {
		if bk.Enabled() == true {
			bk.Write(md)
		}
	}
}

func CloseBackends() {
	mutex.Lock()
	defer mutex.Unlock()

	log.Debugf("backends: ", backends)

	for _, bk := range backends {
		log.Debug("Closing Backend ", bk.Name())
		bk.Close()
	}
}

func RunBackends() {
	mutex.Lock()
	defer mutex.Unlock()

	defer utils.Recover_and_log()

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
	mutex.Lock()
	defer mutex.Unlock()

	backends = nil
}

func CreateBackendsFromRuntimeConf() {
	log.Debug("Create backends from runtime config")
	runtime_conf := config.GetRuntimeConf()
	CreateBackendsFromConf(runtime_conf)
}

func CreateBackendsFromConf(runtime_conf *config.RuntimeConfig) {
	mutex.Lock()
	defer mutex.Unlock()

	defer utils.Recover_and_log()

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
