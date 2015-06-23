package hickwall

import (
	"github.com/oliveagle/hickwall/config"
	"path/filepath"
	"testing"
)

//TODO: how to assert ?
func Test_LoadConfigFromFileAndRun_Single(t *testing.T) {
	config.CONF_FILEPATH, _ = filepath.Abs("./test/config.yml")
	core, _, err := LoadConfigStrategyFile()
	if err != nil {
		t.Error("failed")
		return
	}
	defer core.Close()
}

//TODO: how to assert ?
func Test_LoadConfigFromFileAndRun_GroupDir(t *testing.T) {
	config.CONF_FILEPATH, _ = filepath.Abs("./test/config_wo_groups.yml")
	config.CONF_GROUP_DIRECTORY, _ = filepath.Abs("./test/groups.d")
	core, _, err := LoadConfigStrategyFile()
	if err != nil {
		t.Error("failed")
		return
	}
	if core != nil {
		core.Close()
	}
}

//TODO: how to assert ?
func Test_LoadConfigFromFileAndRun_GroupDir_DupPrefix(t *testing.T) {
	config.CONF_FILEPATH, _ = filepath.Abs("./test/config_wo_groups.yml")
	config.CONF_GROUP_DIRECTORY, _ = filepath.Abs("./test/groups_dup_prefix.d")
	core, _, err := LoadConfigStrategyFile()
	if err == nil {
		t.Error("failed. duplicated group prefix should raise error")
		return
	}
	if core != nil {
		core.Close()
	}
}
