package lib

import (
	"github.com/oliveagle/hickwall/backends"
	// "github.com/oliveagle/hickwall/collectorlib"
	"github.com/oliveagle/hickwall/collectors"
	"github.com/oliveagle/hickwall/config"
	"github.com/oliveagle/hickwall/utils"
	log "github.com/oliveagle/seelog"

	"sync"
	// "time"
)

var (
	mutex sync.Mutex
)

func ReloadWithRuntimeConfig() {
	runtime_conf := config.GetRuntimeConf()
	ReloadConfig(runtime_conf)
}

func ReloadConfig(runtime_conf *config.RuntimeConfig) {
	mutex.Lock()
	defer mutex.Unlock()

	defer utils.Recover_and_log()

	// log.Critical("resp.Config: %+v", runtime_conf)

	log.Info("new config is comming")

	collectors.StopHeartBeat()
	collectors.StopCollectors()
	collectors.RemoveAllCollectors()

	log.Info("Stopped Customized Collectors and Removed them")

	backends.CloseBackends()
	backends.RemoveAllBackends()

	log.Info("Stopped backends and removed them")

	config.UpdateRuntimeConf(runtime_conf)

	log.Info("Updated Runtime Conf with the new one")

	collectors.CreateCollectorsFromRuntimeConf()
	log.Info("Created Customized Colletors")

	backends.CreateBackendsFromRuntimeConf()
	log.Info("Created backends")

	collectors.RunCollectors()
	log.Info("all customized collectors turned on")

	backends.RunBackends()
	log.Info("all backends turned on")

	collectors.StartHeartBeat()
	log.Info("heart beat started")

	log.Info("new config applied")
}

func LoadConfigAndWatching() {
	defer utils.Recover_and_log()

	for resp := range config.WatchConfig() {
		if resp.Err != nil {
			log.Critical("watch config error: %v", resp.Err)
		} else {
			ReloadConfig(resp.Config)
		}
	}
}

// TODO: Load Config Once
