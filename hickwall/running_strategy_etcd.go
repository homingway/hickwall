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

func LoadConfigStrategyEtcd(stop chan error) {
	if stop == nil {
		panic("stop chan is nil")
	}

	respCh := config.WatchRuntimeConfFromEtcd(stop)

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
				core, err := CreateRunningCore(resp.Config)
				// fmt.Println(" -------------- CreateRunningCore finished ------------------------------")
				if err != nil {
					// logging.Errorf("failed to create running core from etcd: %s", err)
					logging.Error(err)
					continue
				} else {
					rconf := resp.Config
					replace_core(core, rconf)
				}
			}
		case <-stop:
			return
		}
	}
}
