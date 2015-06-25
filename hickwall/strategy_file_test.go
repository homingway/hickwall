package hickwall

import (
	"github.com/oliveagle/hickwall/config"
	"path/filepath"
	"testing"
)

func Test_strategy_file_LoadConfigFromFileAndRun_Single(t *testing.T) {
	config.CONF_FILEPATH, _ = filepath.Abs("./test/config.yml")
	_, err := new_core_from_file()
	if err != nil {
		t.Error("failed")
		return
	}
	// panic sucks.
	if the_core == nil {
		t.Error("no core created. ")
		return
	}
	defer close_core()
}

func Test_strategy_file_LoadConfigFromFileAndRun_GroupDir(t *testing.T) {
	config.CONF_FILEPATH, _ = filepath.Abs("./test/config_wo_groups.yml")
	config.CONF_GROUP_DIRECTORY, _ = filepath.Abs("./test/groups.d")
	_, err := new_core_from_file()
	if err != nil {
		t.Error("failed")
		return
	}
	// panic sucks.
	if the_core == nil {
		t.Error("no core created. ")
		return
	}
	defer close_core()
}

func Test_strategy_file_LoadConfigFromFileAndRun_GroupDir_DupPrefix(t *testing.T) {
	config.CONF_FILEPATH, _ = filepath.Abs("./test/config_wo_groups.yml")
	config.CONF_GROUP_DIRECTORY, _ = filepath.Abs("./test/groups_dup_prefix.d")
	_, err := new_core_from_file()
	if err == nil {
		t.Error("failed. duplicated group prefix should raise error")
		return
	}
	if the_core != nil {
		t.Error("should not create a core.")
	}
	defer close_core()
}
