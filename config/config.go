package config

import (
	"fmt"
	"github.com/oliveagle/hickwall/logging"
	"os"
	"path"
	"path/filepath"
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
	_                    = fmt.Sprintln("")
)

type RespConfig struct {
	Config RuntimeConfig
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

	os.MkdirAll(SHARED_DIR, 0755)
	os.MkdirAll(LOG_DIR, 0755)

	logging.InitFileLogger(LOG_FILEPATH[:])

	// try to load core config just ignore the error.
	LoadCoreConfig()
}
