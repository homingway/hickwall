package hickwall

import (
	"github.com/oliveagle/hickwall/config"
	"path/filepath"
	"testing"
)

func Test_strategy_file_LoadConfigFromFileAndRun_Single(t *testing.T) {
	config.CONF_FILEPATH, _ = filepath.Abs("./test/config.yml")
	core, _, err := NewCoreFromFile()
	if err != nil {
		t.Error("failed")
		return
	}
	// panic sucks.
	if core == nil {
		t.Error("no core created. ")
		return
	}
	defer close_core()
}

func Test_strategy_file_LoadConfigFromFileAndRun_GroupDir(t *testing.T) {
	config.CONF_FILEPATH, _ = filepath.Abs("./test/config_wo_groups.yml")
	config.CONF_GROUP_DIRECTORY, _ = filepath.Abs("./test/groups.d")
	core, _, err := NewCoreFromFile()
	if err != nil {
		t.Error("failed")
		return
	}
	// panic sucks.
	if core == nil {
		t.Error("no core created. ")
		return
	}
	defer close_core()
}

func Test_strategy_file_LoadConfigFromFileAndRun_GroupDir_DupPrefix(t *testing.T) {
	config.CONF_FILEPATH, _ = filepath.Abs("./test/config_wo_groups.yml")
	config.CONF_GROUP_DIRECTORY, _ = filepath.Abs("./test/groups_dup_prefix.d")
	core, _, err := NewCoreFromFile()
	if err == nil {
		t.Error("failed. duplicated group prefix should raise error")
		return
	}
	if core != nil {
		t.Error("should not create a core.")
	}
	defer close_core()
}
