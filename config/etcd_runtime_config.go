package config

import (
	//	"fmt"
	"github.com/oliveagle/hickwall/logging"
	"github.com/oliveagle/viper"
	"time"
)

func WatchRuntimeConfFromEtcd(stop chan error) <-chan RespConfig {
	var (
		runtime_viper  = viper.New()
		out            = make(chan RespConfig, 1)
		sleep_duration = time.Second * 5
	)

	if stop == nil {
		stop = make(chan error)
	}

	err := runtime_viper.AddRemoteProvider("etcd", CoreConf.Etcd_url, CoreConf.Etcd_path)
	if err != nil {
		logging.Criticalf("addRemoteProvider Error: %v", err)
	}
	runtime_viper.SetConfigType("YAML")

	go func() {
	label_get_first:
		//need to get config at least once
		var tmp_conf RuntimeConfig
		err = runtime_viper.ReadRemoteConfig()
		if err != nil {
			logging.Errorf("runtime_viper.ReadRemoteConfig Error: %v", err)

			time.Sleep(sleep_duration)
			goto label_get_first
		} else {
			err = runtime_viper.Marshal(&tmp_conf)
			if err != nil {
				logging.Errorf("runtime_viper.Marshal Error: %v", err)

				time.Sleep(sleep_duration)
				goto label_get_first
			} else {
				out <- RespConfig{tmp_conf, nil}
			}
		}

	loop:
		// watch changes
		for {
			var (
				runtime_conf RuntimeConfig
			)

			select {
			case <-stop:
				logging.Debugf("stop watching etcd remote config.")
				break loop
			default:
				logging.Debugf("watching etcd remote config: %s, %s", CoreConf.Etcd_url, CoreConf.Etcd_path)
				err := runtime_viper.WatchRemoteConfig()
				if err != nil {
					logging.Errorf("unable to read remote config: %v", err)
					time.Sleep(sleep_duration)
					continue
				}

				err = runtime_viper.Marshal(&runtime_conf)
				if err != nil {
					logging.Errorf("unable to marshal to config: %v", err)
					time.Sleep(sleep_duration)
					continue
				}

				logging.Debugf("a new config is comming")
				out <- RespConfig{runtime_conf, nil}

				//TODO: make it configurable
				time.Sleep(sleep_duration)
			}
		}
	}()
	return out
}
