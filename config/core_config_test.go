package config

import (
	"path/filepath"
	"testing"
)

func Test_core_config_LoadCoreConfig(t *testing.T) {
	CORE_CONF_FILEPATH, _ = filepath.Abs("./config_example/core_config.yml")
	err := LoadCoreConfig()
	if err != nil {
		t.Error("failed to LoadCoreConfig: %v", err)
		return
	}
	if CoreConf.Rss_limit_mb != 50 || CoreConf.Config_strategy != "file" {
		t.Error("failed to LoadCoreConfig: %+v", CoreConf)
	}

}
