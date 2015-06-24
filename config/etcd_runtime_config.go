//NOTE:  github.com/oliveagle/viper use `hickwall_inuse` branch

package config

import (
	"fmt"
	"github.com/oliveagle/hickwall/logging"
	"github.com/oliveagle/viper"
	_ "github.com/oliveagle/viper/remote"
	"time"
)

var (
	_ = fmt.Sprint("")
)

func WatchRuntimeConfFromEtcd(etcd_url, etcd_path string, stop chan error) <-chan RespConfig {
	var (
		runtime_viper = viper.New()
		out           = make(chan RespConfig, 1)
		// sleep_duration = time.Second * 5
		sleep_duration = time.Second
	)

	if stop == nil {
		stop = make(chan error)
	}

	err := runtime_viper.AddRemoteProvider("etcd", etcd_url, etcd_path)
	if err != nil {
		logging.Criticalf("addRemoteProvider Error: %v", err)
	}
	runtime_viper.SetConfigType("YAML")

	//TODO: should limit retry cnt
	go func() {
		var err error
		var retry_cnt = 0
		var startWatch <-chan time.Time

	label_get_first:
		//need to get config at least once
		var tmp_conf RuntimeConfig

		err = runtime_viper.ReadRemoteConfig()
		if err == nil {
			err = runtime_viper.Marshal(&tmp_conf)
		}

		if err != nil {
			out <- RespConfig{nil, err}
			retry_cnt += 1
			if retry_cnt > 5 {
				out <- RespConfig{nil, fmt.Errorf("cannot get inital config from remote. after 5 attempts")}
				return
			}

			// delay
			time.Sleep(sleep_duration)
			goto label_get_first
		}

		out <- RespConfig{&tmp_conf, nil}

		startWatch = time.Tick(sleep_duration)

	loop:
		// watch changes
		for {
			var runtime_conf RuntimeConfig

			select {
			case <-stop:
				logging.Debugf("stop watching etcd remote config.")
				break loop
			case <-startWatch:
				logging.Debugf("watching etcd remote config: %s, %s", CoreConf.EtcdURL, CoreConf.EtcdPath)
				err := runtime_viper.WatchRemoteConfig()
				if err != nil {
					logging.Errorf("unable to read remote config: %v", err)
					break
				}

				err = runtime_viper.Marshal(&runtime_conf)
				if err != nil {
					logging.Errorf("unable to marshal to config: %v", err)
					break
				}
				logging.Debugf("a new config is comming")
				out <- RespConfig{&runtime_conf, nil}
			}
		}
	}()
	return out
}
