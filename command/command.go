package command

import (
	// "fmt"
	// log "github.com/cihub/seelog"
	// "github.com/codegangsta/cli"
	"github.com/oliveagle/hickwall/config"
	"github.com/oliveagle/hickwall/servicelib"
	// "strings"
)

var PrimaryService = servicelib.NewService("hickwall", config.APP_DESC)

var HelperService = servicelib.NewService("hickwallhelper", "hickwall helper service")
