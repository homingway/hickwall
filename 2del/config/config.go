package config

import (
	"fmt"
	// "github.com/spf13/viper"
	"github.com/oliveagle/viper"
	// "log"
	"os"
	"path"
	"path/filepath"
	// "reflect"

	log "github.com/oliveagle/seelog"

	// "encoding/json"
	// "io/ioutil"
	"time"
)

const (
	VERSION = "v0.0.1"
)

var (
	LOG_DIR            = ""
	SHARED_DIR         = ""
	CONF_FILEPATH      = ""
	CORE_CONF_FILEPATH = ""
	LOG_FILE           = "hickwall.log"
	LOG_FILEPATH       = ""
	CoreConf           CoreConfig // only file
	core_viper         = viper.New()
	rconf              *RuntimeConfig // can retrived from file or etcd
	RuntimeConfChan    = make(chan *RuntimeConfig, 1)
	core_conf_loaded   bool
)

func IsCoreConfigLoaded() bool {
	return CoreConf != nil && core_conf_loaded
}

func LoadCoreConfig() error {
	defer log.Flush()

	initPathes()

	core_viper.SetConfigName("core_config")
	// core_viper.SetConfigFile(CORE_CONF_FILEPATH)
	core_viper.AddConfigPath(SHARED_DIR) // packaged distribution
	core_viper.AddConfigPath(".")        // for hickwall
	core_viper.AddConfigPath("..")       // for hickwall/misc
	core_viper.AddConfigPath("../..")    // for hickwall/misc/try_xxx

	err := core_viper.ReadInConfig()
	if err != nil {
		log.Errorf("No configuration file loaded. core_config.yml :%v", err)
		return fmt.Errorf("No configuration file loaded. core_config.yml :%v", err)
	}

	// log.Debug("core config file used: ", core_viper.ConfigFileUsed())

	err = core_viper.Marshal(&CoreConf)
	if err != nil {
		log.Errorf("Error: unable to parse Core Configuration: %v\n", err)
		return fmt.Errorf("Error: unable to parse Core Configuration: %v\n", err)
	}

	// log.Debug("enable_remote_config: ", CoreConf.Etcd_enabled)
	ConfigLogger()
	if err != nil {
		log.Errorf("LoadCoreConfFile failed: %v", err)
		log.Error("SHARED_DIR: ", SHARED_DIR)
		return fmt.Errorf("LoadCoreConfFile failed: %v", err)

	} else {
		log.Debug("init config, core config loaded")
		log.Debug("LOG_DIR: ", LOG_DIR)
		log.Debug("LOG_FILEPATH: ", LOG_FILEPATH)
	}

	log.Debug("core config file used: ", core_viper.ConfigFileUsed())
	log.Debugf("CoreConfig:  %+v", CoreConf)

	// fmt.Println("core config file used: ", core_viper.ConfigFileUsed())
	// fmt.Println("SHARED_DIR: ", SHARED_DIR)
	core_conf_loaded = true

	log.Debug("CoreConfig Loaded")
	return nil
}

type RespConfig struct {
	Config *RuntimeConfig
	Err    error
}

func loadRuntimeConfFromFile() <-chan *RespConfig {
	log.Debug("loadRuntimeConfFromFile")

	var (
		out           = make(chan *RespConfig, 1)
		runtime_viper = viper.New()
	)
	// runtime_viper.SetConfigFile(config_file)
	runtime_viper.SetConfigName("config")
	runtime_viper.SetConfigType("yaml")
	runtime_viper.AddConfigPath(SHARED_DIR) // packaged distribution
	runtime_viper.AddConfigPath("../..")    // for hickwall/misc/try_xxx
	runtime_viper.AddConfigPath(".")        // for hickwall
	runtime_viper.AddConfigPath("..")       // for hickwall/misc

	go func() {
		var runtime_conf RuntimeConfig

		err := runtime_viper.ReadInConfig()

		log.Debug("RuntimeConfig File Used: ", runtime_viper.ConfigFileUsed())

		// fmt.Println("RuntimeConfig File Used: ", runtime_viper.ConfigFileUsed())

		if err != nil {
			log.Error("loadRuntimeConfFromFile error: ", err)
			out <- &RespConfig{nil, fmt.Errorf("No configuration file loaded. config.yml: %v", err)}
			return
		}

		// Marshal values
		err = runtime_viper.Marshal(&runtime_conf)
		if err != nil {
			log.Error("loadRuntimeConfFromFile error: ", err)
			out <- &RespConfig{nil, fmt.Errorf("Error: unable to parse Configuration: %v\n", err)}
			return
		}

		out <- &RespConfig{&runtime_conf, nil}
		close(out)
		return
	}()

	return out
}

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

func initPathes() {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	dir, _ = filepath.Split(dir)

	SHARED_DIR, _ = filepath.Abs(path.Join(dir, "shared"))
	LOG_DIR, _ = filepath.Abs(path.Join(SHARED_DIR, "logs"))
	LOG_FILEPATH, _ = filepath.Abs(path.Join(LOG_DIR, LOG_FILE))

	CONF_FILEPATH, _ = filepath.Abs(path.Join(SHARED_DIR, "config.yml"))
	CORE_CONF_FILEPATH, _ = filepath.Abs(path.Join(SHARED_DIR, "core_config.yml"))

	// fmt.Println("dir: ", dir)
	// fmt.Println("SHARED_DIR: ", SHARED_DIR)
	// fmt.Println("CONF_FILEPATH: ", CONF_FILEPATH)

	Mkdir_p_logdir(LOG_DIR)
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

func init() {
	LoadCoreConfig()
}

func UpdateRuntimeConf(conf *RuntimeConfig) {
	rconf = conf
}

func GetRuntimeConf() *RuntimeConfig {
	return rconf
}