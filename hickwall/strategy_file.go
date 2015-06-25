package hickwall

import (
	"github.com/oliveagle/hickwall/config"
	"github.com/oliveagle/hickwall/logging"
	//	"github.com/oliveagle/hickwall/newcore"
)

func new_core_from_file() (*config.RuntimeConfig, error) {
	logging.Debug("NewCoreFromFile")
	rconf, err := config.LoadRuntimeConfigFromFiles()
	if err != nil {
		logging.Errorf("NewCoreFromFile: Failed to load RuntimeConfig from files: %v", err)
		return rconf, err
	}
	logging.Debug("NewCoreFromFile: load config from file finished.")
	err = UpdateRunningCore(rconf)
	if err != nil {
		logging.Errorf("NewCoreFromFile: Failed to create running core: %v", err)
		return rconf, err
	}
	logging.Debug("NewCoreFromFile finished witout error")
	return nil, nil
}
