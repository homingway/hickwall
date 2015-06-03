package config

import (
	"fmt"
	"github.com/oliveagle/hickwall/logging"
	"github.com/oliveagle/viper"
)

// public variables ---------------------------------------------------------------
var (
	CoreConf CoreConfig // only file
)

// private variables --------------------------------------------------------------
var (
	core_conf_loaded bool
	core_viper       = viper.New()
)

type CoreConfig struct {
	Log_level         string   `json:"log_level"`         // possible values: trace, debug, info, warn, error, critical
	Log_file_maxsize  int      `json:"log_file_maxsize"`  // TODO: log_file_maxsize
	Log_file_maxrolls int      `json:"log_file_maxrolls"` // TODO: log_file_maxrolls
	Config_Strategy   Strategy `json:"config_strategy"`   // possible values:  file, etcd, registry
	Etcd_url          string   `json:"etcd_url"`          // etcd config
	Etcd_path         string   `json:"etcd_path"`
	Registry_Server   string   `json:"registry_server"` // registry server config
}

func IsCoreConfigLoaded() bool {
	return core_conf_loaded
}

func LoadCoreConfig() error {
	core_viper.SetConfigName("core_config") // core_config.yml
	core_viper.AddConfigPath(SHARED_DIR)    // packaged distribution
	// core_viper.AddConfigPath(".")           // for hickwall
	// core_viper.AddConfigPath("..")          // for hickwall/misc
	// core_viper.AddConfigPath("../..")       // for hickwall/misc/try_xxx

	err := core_viper.ReadInConfig()
	if err != nil {
		logging.Errorf("No configuration file loaded. core_config.yml :%v", err)
		return fmt.Errorf("No configuration file loaded. core_config.yml :%v", err)
	}

	err = core_viper.Marshal(&CoreConf)
	if err != nil {
		logging.Errorf("Error: unable to parse Core Configuration: %v\n", err)
		return fmt.Errorf("Error: unable to parse Core Configuration: %v\n", err)
	}

	logging.SetLevel(CoreConf.Log_level)
	if err != nil {
		logging.Errorf("LoadCoreConfFile failed: %v", err)
		return fmt.Errorf("LoadCoreConfFile failed: %v", err)
	} else {
		logging.Debug("init config, core config loaded")
	}

	logging.Debug("core config file used: ", core_viper.ConfigFileUsed())
	logging.Debugf("CoreConfig:  %+v", CoreConf)

	core_conf_loaded = true

	logging.Debug("CoreConfig Loaded")
	return nil
}
