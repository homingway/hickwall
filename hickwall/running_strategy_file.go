package hickwall

import (
	"github.com/oliveagle/hickwall/config"
	"github.com/oliveagle/hickwall/logging"
	"github.com/oliveagle/hickwall/newcore"
)

func LoadConfigStrategyFile() (newcore.PublicationSet, error) {
	rconf, err := config.LoadRuntimeConfigFromFiles()
	if err != nil {
		logging.Errorf("Failed to load RuntimeConfig from files: %v", err)
		return nil, err
	}
	logging.Info("load config from file finished.")
	core, err := CreateRunningCore(rconf)
	if err != nil {
		logging.Errorf("Failed to create running core: %v", err)
		return nil, err
	}
	logging.Info("replace the core finished.")
	return core, nil
}
