package main

import (
	"fmt"
	"github.com/oliveagle/ole_tryout_daemon/config"
	"github.com/oliveagle/ole_tryout_daemon/servicelib"
	// "github.com/spf13/viper"
	"log"
	"os"
	"strings"
)

const (
	version = "v0.0.1"
	svcName = "oleservice"
	svcDesc = "ole service description"
	port    = ":9977"
)

func usage(errmsg string) {
	fmt.Fprintf(os.Stderr,
		"%s\n\n"+
			"usage: %s <command>\n"+
			"       where <command> is one of\n"+
			"       install, remove, status, start, stop, pause or continue.\n",
		errmsg, os.Args[0])
	os.Exit(2)
}

func main() {
	// - log --------------------
	f, err := os.OpenFile(getLogFilePath(), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v \n", err)
	}
	defer f.Close()
	log.SetOutput(f)
	// -------------------- log -

	config.SetDefault()

	srv := servicelib.NewService(svcName, svcDesc)

	if len(os.Args) >= 2 {
		log.Println("new func main\r\n")

		cmd := strings.ToLower(os.Args[1])
		switch cmd {
		case "install":
			err = srv.InstallService()
		case "remove":
			err = srv.RemoveService()
		case "start":
			err = srv.StartService()
		case "stop":
			err = srv.StopService()
		case "pause":
			err = srv.PauseService()
		case "continue":
			err = srv.ContinueService()
		case "status":
			err = srv.Status()
		case "config":
			err = srv.Config()
		default:
			usage(fmt.Sprintf("invalid command %s", cmd))
		}
		if err != nil {
			log.Fatalf("failed to %s %s: %v", cmd, svcName, err)
		}
	} else {
		isIntSess, err := srv.IsAnInteractiveSession()
		if err != nil {
			log.Fatalf("failed to determine if we are running in an interactive session: %v", err)
		}
		if !isIntSess {
			runService(svcName, false)
			return
		}
		// runService(svcName, false)
	}

	return
}
