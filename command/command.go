package command

import (
	"github.com/oliveagle/hickwall/servicelib"
)

var PrimaryService = servicelib.NewService("hickwall", "monitoring system")
