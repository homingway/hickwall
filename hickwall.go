package main

import (
	"fmt"
	log "github.com/oliveagle/hickwall/_third_party/seelog"
	"github.com/oliveagle/hickwall/config"
	"github.com/oliveagle/hickwall/servicelib"
	"os"
	"strings"
)

func usage(errmsg string) {
	fmt.Fprintf(os.Stderr,
		"%s\n\n"+
			"usage: %s <command>\n"+
			"       where <command> is one of\n"+
			"       install, remove, status, start, stop, pause or continue.\n",
		errmsg, os.Args[0])
}

var err error

func main() {
	defer log.Flush()

	log.Error("error")

	srv := servicelib.NewService(config.APP_NAME, config.APP_DESC)

	if len(os.Args) >= 2 {
		cmd := strings.ToLower(os.Args[1])
		err = servicelib.HandleCmd(srv, cmd)
		if err != nil {
			usage(fmt.Sprintf("failed to %s %s: %v", cmd, config.APP_NAME, err))
			return
		}
	} else {
		isIntSess, err := srv.IsAnInteractiveSession()
		if err != nil {
			log.Error("failed to determine if we are running in an interactive session: %v", err)
			return
		}
		if !isIntSess {
			runService(config.APP_NAME, false)
			return
		}
		// runService(config.APP_NAME, false)
	}

	return
}
