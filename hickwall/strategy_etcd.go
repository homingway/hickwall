package hickwall

import (
	"fmt"
	"github.com/kr/pretty"
	"github.com/oliveagle/hickwall/config"
	"github.com/oliveagle/hickwall/logging"
)

var (
	_ = fmt.Sprint("")
)

func NewCoreFromEtcd(etcd_machines []string, etcd_path string, stop chan error) {
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

	// loop:
	for {
		select {
		case resp := <-respCh:
			// fmt.Println(" ------------ got resp --------------------")
			pretty.Println(resp)
			if resp.Err != nil {
				logging.Error(resp.Err)
				// logging.Errorf("failed to get RuntimeConf from etcd: %v", resp.Err)
				continue
			} else {
				close_core() //TODO: to prevent race condition. maybe we can safely remove this line.
				_, err := UpdateRunningCore(resp.Config)
				// fmt.Println(" -------------- CreateRunningCore finished ------------------------------")
				if err != nil {
					// logging.Errorf("failed to create running core from etcd: %s", err)
					logging.Error(err)
					continue
				} else {
					rconf := resp.Config
					//					replace_core(core, rconf)

					// dump cached runtime config only if it changed.
					if cached_hash != rconf.GetHash() {
						err = config.DumpRuntimeConfig(rconf)
						if err != nil {
							logging.Errorf("failed to dump runtime config: %v", err)
						}
						cached_hash = rconf.GetHash()
					}

				}
			}
		case <-stop:
			return
		}
	}
}
