package main

import (
	"github.com/oliveagle/hickwall/backends"
	"github.com/oliveagle/hickwall/collectorlib"
	"github.com/oliveagle/hickwall/collectors"
	"github.com/oliveagle/hickwall/config"
	"github.com/oliveagle/hickwall/utils"
	log "github.com/oliveagle/seelog"
)

func LoadConfigAndReload(mdCh chan<- collectorlib.MultiDataPoint) {
	defer utils.Recover_and_log()

	for resp := range config.WatchConfig() {
		if resp.Err != nil {
			log.Critical("watch config error: %v", resp.Err)
		} else {
			defer log.Flush()

			log.Debug("new config is comming")

			collectors.StopCustomizedCollectors()
			collectors.RemoveAllCustomizedCollectors()

			log.Debug("Stopped Customized Collectors and Removed them")

			backends.CloseBackends()
			backends.RemoveAllBackends()

			log.Debug("Stopped backends and removed them")

			config.UpdateRuntimeConf(resp.Config)

			log.Debug("Updated Runtime Conf with the new one")

			collectors.CreateCustomizedCollectorsFromRuntimeConf()
			log.Debug("Created Customized Colletors")

			backends.CreateBackendsFromRuntimeConf()
			log.Debug("Created backends")

			collectors.RunCustomizedCollectors(mdCh)
			log.Debug("all customized collectors turned on")

			backends.RunBackends()
			log.Debug("all backends turned on")

			log.Debug("new config applied")
		}
	}
}
