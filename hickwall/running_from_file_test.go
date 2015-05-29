package hickwall

import (
	"github.com/oliveagle/hickwall/config"
	"path/filepath"
	"testing"
)

func Test_LoadConfigFromFileAndRun(t *testing.T) {
	config.CONF_FILEPATH, _ = filepath.Abs("./test/config.yml")
	core, err := LoadConfigFromFileAndRun()
	if err != nil {
		t.Error("failed")
		return
	}
	defer core.Close()
}
