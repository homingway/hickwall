package config

import (
	"fmt"
	"github.com/oliveagle/hickwall/logging"
	"github.com/oliveagle/viper"
	"os"
)

// public variables ---------------------------------------------------------------
var (
	CoreConf CoreConfig // only file
	_        = fmt.Sprintln("")
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

	file, err := os.Open(CORE_CONF_FILEPATH)
	if err != nil {
		return fmt.Errorf("cannot open file file: %s - err: %v", CORE_CONF_FILEPATH, err)
	}
	defer file.Close()

	err = core_viper.ReadConfig(file)
	if err != nil {
		return fmt.Errorf("No configuration file loaded. core_config.yml :%v", err)
	}

	err = core_viper.Marshal(&CoreConf)
	if err != nil {
		return fmt.Errorf("Error: unable to parse Core Configuration: %v\n", err)
	}

	logging.SetLevel(CoreConf.Log_level)
	if err != nil {
		return fmt.Errorf("LoadCoreConfFile failed: %v", err)
	}

	logging.Debugf("core config file used: %s\n", core_viper.ConfigFileUsed())
	logging.Debugf("SHARED_DIR:            %s\n", SHARED_DIR)
	logging.Debugf("LOG_DIR:               %s\n", LOG_DIR)
	logging.Debugf("LOG_FILEPATH:          %s\n", LOG_FILEPATH)
	logging.Debugf("CORE_CONF_FILEPATH:    %s\n", CORE_CONF_FILEPATH)
	logging.Debugf("CONF_FILEPATH:         %s\n", CONF_FILEPATH)
	logging.Debugf("REGISTRY_FILEPATH:     %s\n", REGISTRY_FILEPATH)
	logging.Debugf("CONF_GROUP_DIRECTORY:  %s\n", CONF_GROUP_DIRECTORY)
	logging.Debugf("CoreConfig:            %+v\n", CoreConf)
	logging.Debug("CoreConfig Loaded ==============================================")

	core_conf_loaded = true
	return nil
}
