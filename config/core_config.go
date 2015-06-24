package config

import (
	"fmt"
	"github.com/oliveagle/hickwall/logging"
	"github.com/oliveagle/hickwall/newcore"
	"github.com/oliveagle/hickwall/utils"
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
	Hostname        string   `json:"hostname"`          // hostname, which will override auto collected hostname
	RssLimitMb      int      `json:"rss_limit_mb"`      // rss_limit_mb to kill client, default is 50Mb
	LogLevel        string   `json:"log_level"`         // possible values: trace, debug, info, warn, error, critical
	LogFileMaxSize  int      `json:"log_file_maxsize"`  // TODO: log_file_maxsize
	LogFileMaxRolls int      `json:"log_file_maxrolls"` // TODO: log_file_maxrolls
	ConfigStrategy  Strategy `json:"config_strategy"`   // possible values:  file, etcd, registry
	EtcdURL         string   `json:"etcd_url"`          // etcd url
	EtcdPath        string   `json:"etcd_path"`         // etcd config path
	RegistryURL     string   `json:"registry_url"`      // registry server config
	ListenPort      int      `json:"listen_port"`       // api listen port, default 3031
	SecureAPIWrite  bool     `json:"secure_api_write"`  // default false, use admin server public key to protect write apis.
	SecureAPIRead   bool     `json:"secure_api_read"`   // default false, use admin server public key to protect read apis.
	ServerPubKey    string   `json:"server_pub_key"`    // use this public key to verify server api call
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
		// fmt.Printf("cannot open file file: %s - err: %v\n", CORE_CONF_FILEPATH, err)
		return fmt.Errorf("cannot open file file: %s - err: %v", CORE_CONF_FILEPATH, err)
	}
	defer file.Close()

	err = core_viper.ReadConfig(file)
	if err != nil {
		// fmt.Printf("No configuration file loaded. core_config.yml :%v\n", err)
		return fmt.Errorf("No configuration file loaded. core_config.yml :%v", err)
	}

	err = core_viper.Marshal(&CoreConf)
	if err != nil {
		// fmt.Printf("Error: unable to parse Core Configuration: %v\n", err)
		return fmt.Errorf("Error: unable to parse Core Configuration: %v\n", err)
	}

	if CoreConf.RssLimitMb <= 0 {
		CoreConf.RssLimitMb = 50 //deffault rss limit
	}
	if CoreConf.ListenPort <= 0 {
		CoreConf.ListenPort = 3031
	}
	if CoreConf.Hostname != "" {
		newcore.SetHostname(CoreConf.Hostname)
	}

	logging.SetLevel(CoreConf.LogLevel)
	if err != nil {
		return fmt.Errorf("LoadCoreConfFile failed: %v", err)
	}

	if CoreConf.SecureAPIRead || CoreConf.SecureAPIWrite {
		// we should check public key config.
		_, err := utils.LoadPublicKeyFromPath(CoreConf.ServerPubKey)
		if err != nil {
			logging.Criticalf("unable to load server public key while SecureAPIx is set to be true: %s", err)
		}
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
