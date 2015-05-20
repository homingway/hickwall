package command

import (
	"github.com/oliveagle/hickwall/config"
	"github.com/oliveagle/hickwall/servicelib"
	"path/filepath"
)

var HelperService = servicelib.NewServiceFromPath("hickwallhelper", "hickwall helper service", filepath.Join(config.SHARED_DIR, "hickwall_helper.exe"))
