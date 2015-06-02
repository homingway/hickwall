package hickwall

import (
	"github.com/oliveagle/hickwall/config"
	"github.com/oliveagle/hickwall/logging"
)

func LoadConfigStrategyEtcd(stop chan chan error) {
	stopCh := make(chan error)
	respCh := config.WatchRuntimeConfFromEtcd(stopCh)

	// loop:
	for {
		select {
		case resp := <-respCh:
			if resp.Err != nil {
				logging.Error("failed to get RuntimeConf from etcd: ", resp.Err)
				continue
			} else {
				close_core() //TODO: to prevent race condition. maybe we can safely remove this line.
				core, err := CreateRunningCore(&resp.Config)
				if err != nil {
					logging.Error("failed to create running core from etcd: ", err)
					continue
				} else {
					replace_core(core)
				}
			}
		case errc := <-stop:
			errc <- nil
			return
		}
	}
}
