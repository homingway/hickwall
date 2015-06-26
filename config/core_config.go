package config

import (
	"fmt"
	"github.com/oliveagle/hickwall/logging"
	"github.com/oliveagle/hickwall/newcore"
	"github.com/oliveagle/hickwall/utils"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// public variables ---------------------------------------------------------------
var (
	CoreConf CoreConfig // only file
	_        = fmt.Sprintln("")
)

// private variables --------------------------------------------------------------
var (
	core_conf_loaded bool
)

type CoreConfig struct {
	Hostname            string   `json:"hostname"`            // hostname, which will override auto collected hostname
	Rss_limit_mb        int      `json:"rss_limit_mb"`        // rss_limit_mb to kill client, default is 50Mb
	Log_level           string   `json:"log_level"`           // possible values: trace, debug, info, warn, error, critical
	Log_file_maxsize    int      `json:"log_file_maxsize"`    // TODO: log_file_maxsize
	Log_file_maxrolls   int      `json:"log_file_maxrolls"`   // TODO: log_file_maxrolls
	Config_strategy     string   `json:"config_strategy"`     // possible values:  file, etcd, registry
	Etcd_machines       []string `json:"etcd_machines"`       // etcd machines
	Etcd_path           string   `json:"etcd_path"`           // etcd config path
	Registry_urls       []string `json:"registry_urls"`       // registry server config
	Enable_http_api     bool     `json:"enable_http_api"`     // enable http api. default is false.
	Listen_port         int      `json:"listen_port"`         // api listen port, default 3031
	Secure_api_write    bool     `json:"secure_api_write"`    // default false, use admin server public key to protect write apis.
	Secure_api_read     bool     `json:"secure_api_read"`     // default false, use admin server public key to protect read apis.
	Server_pub_key_path string   `json:"server_pub_key_path"` // use this public key to verify server api call
}

func IsCoreConfigLoaded() bool {
	return core_conf_loaded
}

func LoadCoreConfig() error {
	data, err := ioutil.ReadFile(CORE_CONF_FILEPATH)
	if err != nil {
		return fmt.Errorf("faild to read core config: %v", err)
	}
	CoreConf = CoreConfig{}
	err = yaml.Unmarshal(data, &CoreConf)
	if err != nil {
		return fmt.Errorf("unable to unmarshal yaml: %v", err)
	}

	if CoreConf.Rss_limit_mb <= 0 {
		CoreConf.Rss_limit_mb = 50 //deffault rss limit
	}
	if CoreConf.Listen_port <= 0 {
		CoreConf.Listen_port = 3031
	}
	if CoreConf.Hostname != "" {
		newcore.SetHostname(CoreConf.Hostname)
	}

	logging.SetLevel(CoreConf.Log_level)
	if err != nil {
		return fmt.Errorf("LoadCoreConfFile failed: %v", err)
	}

	if CoreConf.Enable_http_api && (CoreConf.Secure_api_read || CoreConf.Secure_api_write) {
		// we should check public key config.
		_, err := utils.LoadPublicKeyFromPath(CoreConf.Server_pub_key_path)
		if err != nil {
			logging.Criticalf("unable to load server public key while SecureAPIx is set to be true: %s", err)
		}
	}

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
