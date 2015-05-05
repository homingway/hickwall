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
	config_path        []string
	LOG_DIR            = ""
	SHARED_DIR         = ""
	CONF_FILEPATH      = ""
	CORE_CONF_FILEPATH = ""

	LOG_FILE     = "hickwall.log"
	LOG_FILEPATH = ""

	CoreConf   CoreConfig // only file
	core_viper = viper.New()

	rconf Config // can retrived from file or etcd

	RuntimeConfChan = make(chan *Config, 1)
)

func loadCoreConfig() {

	initPathes()

	core_viper.SetConfigName("core_config")
	// core_viper.SetConfigFile(CORE_CONF_FILEPATH)
	core_viper.AddConfigPath(SHARED_DIR) // packaged distribution
	core_viper.AddConfigPath(".")        // for hickwall
	core_viper.AddConfigPath("..")       // for hickwall/misc
	core_viper.AddConfigPath("../..")    // for hickwall/misc/try_xxx

	// err := LoadCoreConfig()
	err := core_viper.ReadInConfig()
	if err != nil {
		log.Errorf("No configuration file loaded. core_config.yml :%v", err)
		log.Flush()
		os.Exit(1)
	}

	// log.Debug("core config file used: ", core_viper.ConfigFileUsed())

	err = core_viper.Marshal(&CoreConf)
	if err != nil {
		log.Errorf("Error: unable to parse Core Configuration: %v\n", err)
		log.Flush()
		os.Exit(1)
	}

	// log.Debug("enable_remote_config: ", CoreConf.Etcd_enabled)

	ConfigLogger()
	if err != nil {
		log.Errorf("LoadCoreConfFile failed: %v", err)
		log.Error("SHARED_DIR: ", SHARED_DIR)
		// log.Error("CORE_CONF_FILEPATH: ", CORE_CONF_FILEPATH)
		log.Flush()
		os.Exit(1)
	} else {
		log.Debug("init config, core config loaded")
		log.Debug("LOG_DIR: ", LOG_DIR)
		log.Debug("LOG_FILEPATH: ", LOG_FILEPATH)
	}

	log.Debug("core config file used: ", core_viper.ConfigFileUsed())
	log.Debugf("CoreConfig:  %+v", CoreConf)
}

type RespConfig struct {
	Config *Config
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
		var runtime_conf Config

		err := runtime_viper.ReadInConfig()

		log.Debug("Runtime Config File Used: ", runtime_viper.ConfigFileUsed())
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
		var tmp_conf Config
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

		// watch changes
		for {
			var (
				runtime_conf Config
			)

			select {
			case <-stop:
				log.Debugf("stop watching etcd remote config.")
				break
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
	// SHARED_DIR = `D:\Users\rhtang\oledev\gocodez\src\github.com\oliveagle\shared`

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

func init() {
	loadCoreConfig()
	log.Debug("CoreConfig Loaded")
}

func UpdateRuntimeConf(conf *Config) {
	rconf = *conf
}

func GetRuntimeConf() *Config {
	return &rconf
}
