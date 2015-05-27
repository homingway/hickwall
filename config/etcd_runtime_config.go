package config

import (
	//	"fmt"
	log "github.com/oliveagle/seelog"
	"github.com/oliveagle/viper"
	"time"
)

func WatchRuntimeConfFromEtcd(stop chan bool) <-chan *RespConfig {
	var (
		runtime_viper = viper.New()
		out           = make(chan *RespConfig, 1)
	)

	if stop == nil {
		stop = make(chan bool)
	}

	err := runtime_viper.AddRemoteProvider("etcd", CoreConf.Etcd_url, CoreConf.Etcd_path)
	if err != nil {
		log.Criticalf("addRemoteProvider Error: %v", err)
	}
	runtime_viper.SetConfigType("YAML")

	go func() {
		//need to get config at least once
		var tmp_conf RuntimeConfig
		err = runtime_viper.ReadRemoteConfig()
		if err != nil {
			log.Errorf("runtime_viper.ReadRemoteConfig Error: %v", err)
		} else {
			err = runtime_viper.Marshal(&tmp_conf)
			if err != nil {
				log.Errorf("runtime_viper.Marshal Error: %v", err)
			} else {
				out <- &RespConfig{&tmp_conf, nil}
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
				log.Debugf("stop watching etcd remote config.")
				break loop
			default:
				log.Debugf("watching etcd remote config: %s, %s", CoreConf.Etcd_url, CoreConf.Etcd_path)
				err := runtime_viper.WatchRemoteConfig()
				if err != nil {
					log.Errorf("unable to read remote config: %v", err)
					time.Sleep(time.Second * 5)
					continue
				}

				err = runtime_viper.Marshal(&runtime_conf)
				if err != nil {
					log.Errorf("unable to marshal to config: %v", err)
					time.Sleep(time.Second * 5)
					continue
				}

				log.Debugf("a new config is comming")
				out <- &RespConfig{&runtime_conf, nil}

				//TODO: make it configurable
				time.Sleep(time.Second * 5)
			}
		}
	}()
	return out
}
