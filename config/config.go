package config

import (
	//	"fmt"
	// "github.com/spf13/viper"
	//	"github.com/oliveagle/viper"
	// "log"
	"os"
	"path"
	"path/filepath"
	// "reflect"

	//	log "github.com/oliveagle/seelog"

	// "encoding/json"
	// "io/ioutil"
	//	"time"
)

const (
	VERSION = "v0.0.1"
)

var (
	LOG_DIR              = ""
	LOG_FILE             = "hickwall.log"
	LOG_FILEPATH         = ""
	SHARED_DIR           = ""
	CORE_CONF_FILEPATH   = ""
	CONF_FILEPATH        = ""
	CONF_GROUP_DIRECTORY = ""
	REGISTRY_FILEPATH    = ""
)

type RespConfig struct {
	Config *RuntimeConfig
	Err    error
}

func init() {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	dir, _ = filepath.Split(dir)

	SHARED_DIR, _ = filepath.Abs(path.Join(dir, "shared"))
	LOG_DIR, _ = filepath.Abs(path.Join(SHARED_DIR, "logs"))
	LOG_FILEPATH, _ = filepath.Abs(path.Join(LOG_DIR, LOG_FILE))
	CORE_CONF_FILEPATH, _ = filepath.Abs(path.Join(SHARED_DIR, "core_config.yml"))
	CONF_FILEPATH, _ = filepath.Abs(path.Join(SHARED_DIR, "config.yml"))
	REGISTRY_FILEPATH, _ = filepath.Abs(path.Join(SHARED_DIR, "registry"))

	// CollectorConfigGroup with in this folder
	CONF_GROUP_DIRECTORY, _ = filepath.Abs(path.Join(SHARED_DIR, "groups.d"))

	// fmt.Println("dir: ", dir)
	// fmt.Println("SHARED_DIR: ", SHARED_DIR)
	// fmt.Println("CONF_FILEPATH: ", CONF_FILEPATH)
	// fmt.Println("CONF_DIRECTORY: ", CONF_DIRECTORY)

	Mkdir_p_logdir(LOG_DIR)

	// we don't need to always load core config
	//	LoadCoreConfig()
}

//func WatchConfig() <-chan *RespConfig {
//	if CoreConf.Etcd_enabled == true {
//		return WatchRuntimeConfFromEtcd(nil)
//	} else {
//		return loadRuntimeConfFromFile()
//	}
//}
