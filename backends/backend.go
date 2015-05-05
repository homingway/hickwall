package backends

import (
	"fmt"
	"github.com/oliveagle/hickwall/collectorlib"
	"github.com/oliveagle/hickwall/config"
	log "github.com/oliveagle/seelog"
	"strings"
)

var (
	// backends = make(map[string]TSWriter)

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
// 	backends = map[string]TSWriter{}
// }

// 	var runtime_conf = config.GetRuntimeConf()

// 	backends["stdout"] = NewStdoutWriter(runtime_conf.Transport_stdout)

// 	backends["file"] = NewFileWriter(runtime_conf.Transport_file)

// 	// influxdb backends
// 	for _, iconf := range runtime_conf.Transport_influxdb {
// 		bkname := fmt.Sprintf("influxdb-%s", influxdbParseVersionFromString(iconf.Version))
// 		bk, err := NewInfluxdbWriter(iconf)
// 		if err != nil {
// 			log.Criticalf("create backend failed: %v ", err)
// 			continue
// 		}
// 		backends[bkname] = bk
// 	}

// }

func GetBackendList() (res []string) {
	for _, bk := range backends {
		res = append(res, bk.Name())
	}
	// for key, _ := range backends {
	// 	res = append(res, key)
	// }
	return
}

func GetBackendByName(name string) (w TSWriter, b bool) {
	for _, bk := range backends {
		if strings.ToLower(bk.Name()) == strings.ToLower(name) {
			return bk, true
		}
	}
	// w, b = backends[strings.ToLower(name)]
	return nil, false
}

func GetBackendByNameVersion(name, version string) (TSWriter, bool) {
	key := strings.Join([]string{name, version}, "-")
	// w, b = backends[strings.ToLower(key)]
	return GetBackendByName(strings.ToLower(key))
}

func WriteToBackends(md collectorlib.MultiDataPoint) {
	for _, bk := range backends {
		if bk.Enabled() == true {
			log.Debug("Backend.Write.Endabled: ", bk.Name())
			bk.Write(md)
		}
	}

	// for key, backend := range backends {
	// 	if backend.Enabled() == true {
	// 		log.Debug("Backend.Write.Endabled: ", key)
	// 		backend.Write(md)
	// 	}
	// 	// else {
	// 	// 	log.Debug("Backend.Write.Disabled: ", key)
	// 	// }
	// }
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

func CreateBackendsFromConf(runtime_conf *config.Config) {
	defer log.Flush()

	log.Debug("create backends from conf: %v", runtime_conf)
	log.Flush()

	// backends["stdout"] = NewStdoutWriter(runtime_conf.Transport_stdout)
	backends = append(backends, NewStdoutWriter("stdout", runtime_conf.Transport_stdout))
	log.Debug("stdout backend created")

	// backends["file"] = NewFileWriter(runtime_conf.Transport_file)
	backends = append(backends, NewFileWriter("file", runtime_conf.Transport_file))
	log.Debug("file backend created")

	// influxdb backends
	for _, iconf := range runtime_conf.Transport_influxdb {
		bkname := fmt.Sprintf("influxdb-%s", influxdbParseVersionFromString(iconf.Version))
		log.Debugf("Creating backend: %s", bkname)
		bk, err := NewInfluxdbWriter(bkname, iconf)
		if err != nil {
			log.Criticalf("create backend failed: %v ", err)
			log.Flush()
			continue
		}
		log.Debugf("backend created: %s", bkname)
		log.Flush()
		backends = append(backends, bk)
		log.Debugf("backend add to map: %s", bkname)
		log.Flush()

	}

	fmt.Println("hh")
	log.Debug("all backends created")
	log.Flush()
}
