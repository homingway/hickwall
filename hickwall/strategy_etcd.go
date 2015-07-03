package hickwall

import (
	"fmt"
	"github.com/oliveagle/hickwall/config"
	"github.com/oliveagle/hickwall/logging"
)

var (
	_ = fmt.Sprint("")
)

func new_core_from_etcd(etcd_machines []string, etcd_path string, stop chan error) {
	logging.Debug("new_core_from_etcd started")
	if stop == nil {
		panic("stop chan is nil")
	}
	if len(etcd_machines) <= 0 {
		logging.Critical("EtcdMachines is empty!!")
		panic("LoadConfigStrategyEtcd: EtcdMachines is empty!!")
	}
	if etcd_path == "" {
		logging.Critical("EtcdPath is empty!!")
		panic("LoadConfigStrategyEtcd: EtcdPath is empty!!")
	}

	cached_hash, _ := config.GetCachedRuntimeConfigHash()

	respCh := config.WatchRuntimeConfFromEtcd(etcd_machines, etcd_path, stop)

	for {
		select {
		case resp := <-respCh:
			logging.Info("NewCoreFromEtcd: a new response is comming.")
			if resp.Err != nil {
				logging.Error(resp.Err)
				break
			} else {
				err := UpdateRunningCore(resp.Config)
				if err != nil {
					logging.Error(err)
					break
				} else {
					// dump cached runtime config only if it changed and working
					if cached_hash != resp.Config.GetHash() {
						err = config.DumpRuntimeConfig(resp.Config)
						if err != nil {
							logging.Errorf("failed to dump runtime config: %v", err)
						}
						cached_hash = resp.Config.GetHash()
					}
				}
				logging.Info("NewCoreFromEtcd: replaced new core and updated cache.")
			}
		case <-stop:
			return
		}
	}
}
