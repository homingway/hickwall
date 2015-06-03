package hickwall

import (
	"fmt"
	//	"github.com/oliveagle/hickwall/collectors"
	"github.com/oliveagle/hickwall/config"
	"github.com/oliveagle/hickwall/newcore"
)

func LoadConfigStrategyFile() (newcore.PublicationSet, error) {
	rconf, err := config.LoadRuntimeConfigFromFiles()
	if err != nil {
		fmt.Println("Failed to load RuntimeConfig from files: ", err)
		return nil, err
	}
	core, err := CreateRunningCore(rconf)
	if err != nil {
		fmt.Println("Failed to create running core: ", err)
		return nil, err
	}

	return core, nil
}
