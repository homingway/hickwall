package config

func WatchRuntimeConfFromEtcd(stop chan bool) <-chan *RespConfig {
	var (
		runtime_viper = viper.New()
		out           = make(chan *RespConfig, 1)
	)

	if stop == nil {
		stop = make(chan bool)
	}

	err := runtime_viper.AddRemoteProvider("etcd", CoreConf.Etcd_url, CoreConf.Etcd_path)
	// err := runtime_viper.AddRemoteProvider("etcd", "http://192.168.59.103:4001", "/config/host/DST54869.yml")
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

func WatchConfig() <-chan *RespConfig {
	if CoreConf.Etcd_enabled == true {
		return WatchRuntimeConfFromEtcd(nil)
	} else {
		return loadRuntimeConfFromFile()
	}
}

func LoadRuntimeConfFromFileOnce() error {
	defer log.Flush()

	for resp := range loadRuntimeConfFromFile() {
		if resp.Err != nil {
			log.Errorf("cannot load runtime config from file: %v", resp.Err)
			return fmt.Errorf("cannot load runtime config from file: %v", resp.Err)
		} else {
			UpdateRuntimeConf(resp.Config)
			log.Debug("updated runtime config")
		}
	}
	return nil
}
