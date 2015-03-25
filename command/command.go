package command

import (
	"github.com/oliveagle/hickwall/servicelib"
)

var PrimaryService = servicelib.NewService("hickwall", "monitoring system")

var HelperService = servicelib.NewService("hickwallhelper", "hickwall helper service")
