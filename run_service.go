package main

import (
	"github.com/oliveagle/hickwall/backends"
	// "github.com/oliveagle/hickwall/collectorlib"
	"github.com/oliveagle/hickwall/collectors"
	"github.com/oliveagle/hickwall/config"
	"github.com/oliveagle/hickwall/utils"
	log "github.com/oliveagle/seelog"

	"sync"
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

	log.Critical("resp.Config: %+v", runtime_conf)

	log.Debug("new config is comming")

	collectors.StopCollectors()
	collectors.RemoveAllCollectors()

	log.Debug("Stopped Customized Collectors and Removed them")

	backends.CloseBackends()
	backends.RemoveAllBackends()

	log.Debug("Stopped backends and removed them")

	config.UpdateRuntimeConf(runtime_conf)

	log.Debug("Updated Runtime Conf with the new one")

	collectors.CreateCollectorsFromRuntimeConf()
	log.Debug("Created Customized Colletors")

	backends.CreateBackendsFromRuntimeConf()
	log.Debug("Created backends")

	collectors.RunCollectors()
	log.Debug("all customized collectors turned on")

	backends.RunBackends()
	log.Debug("all backends turned on")

	log.Debug("new config applied")
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
