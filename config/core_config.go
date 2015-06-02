package config

import (
	"fmt"
	"github.com/oliveagle/hickwall/logging"
	"github.com/oliveagle/viper"
	"strings"
)

// DataType is the type of rate for a metric: gauge, counter, or rate.
type Strategy string

const (
	FILE     Strategy = "file"
	ETCD              = "etcd"
	REGISTRY          = "registry"
)

func (s *Strategy) IsValid() bool {
	k := strings.ToLower(string(*s))
	switch k {
	case "file", "etcd", "registry":
		return true
	default:
		return false
	}
}

func (s *Strategy) GetString() string {
	return strings.ToLower(string(*s))
}

// public variables
var (
	CoreConf CoreConfig // only file
)

// private variables
var (
	core_conf_loaded bool
	core_viper       = viper.New()
)

type CoreConfig struct {
	Log_level         string `json:"log_level"`
	Log_file_maxsize  int    `json:"log_file_maxsize"`
	Log_file_maxrolls int    `json:"log_file_maxrolls"`

	Heart_beat_interval string `json:"heart_beat_interval"`

	// possible values:  file, etcd, registry
	Config_Strategy Strategy `json:"config_strategy"`

	// etcd config
	Etcd_url  string `json:"etcd_url"`
	Etcd_path string `json:"etcd_path"`

	// registry server config
	Registry_Server string `json:"registry_server"`
}

func IsCoreConfigLoaded() bool {
	return core_conf_loaded
}

func LoadCoreConfig() error {

	// defer logging.Flush()

	core_viper.SetConfigName("core_config")
	// core_viper.SetConfigFile(CORE_CONF_FILEPATH)
	core_viper.AddConfigPath(SHARED_DIR) // packaged distribution
	core_viper.AddConfigPath(".")        // for hickwall
	core_viper.AddConfigPath("..")       // for hickwall/misc
	core_viper.AddConfigPath("../..")    // for hickwall/misc/try_xxx

	err := core_viper.ReadInConfig()
	if err != nil {
		logging.Errorf("No configuration file loaded. core_config.yml :%v", err)
		return fmt.Errorf("No configuration file loaded. core_config.yml :%v", err)
	}

	// logging.Debug("core config file used: ", core_viper.ConfigFileUsed())

	err = core_viper.Marshal(&CoreConf)
	if err != nil {
		logging.Errorf("Error: unable to parse Core Configuration: %v\n", err)
		return fmt.Errorf("Error: unable to parse Core Configuration: %v\n", err)
	}

	ConfigLogger()
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
