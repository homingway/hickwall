package logging

import (
	// "github.com/oliveagle/hickwall/logging/config"
	"github.com/op/go-logging"
	"os"
	"path"
	"runtime"
	"strings"
)

var format = logging.MustStringFormatter(
	"%{color}%{time} %{shortfunc} > %{level:.4s} %{id:03x}%{color:reset} %{message}",
)

func MustGetLogger() *logging.Logger {
	return logging.MustGetLogger(GetModule())
}

func GetModule() string {
	_, filename, _, _ := runtime.Caller(1)
	res := path.Dir(filename)
	idx := strings.LastIndex(res, "github.com")
	if idx > 0 {
		res = res[idx:]
	}
	return res
}

func init() {
	backend1 := logging.NewLogBackend(os.Stderr, "", 0)
	// backend2 := logging.NewLogBackend(os.Stderr, "", 0)

	leveled1 := logging.AddModuleLevel(backend1)
	// leveled2 := logging.AddModuleLevel(backend2)

	multi := logging.MultiLogger(leveled1)
	// multi := logging.MultiLogger(leveled1, leveled2)

	multi.SetLevel(logging.ERROR, "github.com/oliveagle/hickwall/logging")

	logging.SetBackend(multi)
	logging.SetFormatter(format)
}
