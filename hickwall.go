package main

import (
	"fmt"
	"github.com/oliveagle/hickwall/config"
	"github.com/oliveagle/hickwall/servicelib"
	// l4g "github.com/oliveagle/log4go"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
	"strings"
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

func mkdir_p_logdir(logfile string) {
	dir, _ := filepath.Split(logfile)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		fmt.Println("Error: cannot create log dir: %s, err: %s", dir, err)
	}
}

func main() {
	config.LoadConfig()
	logfile := viper.GetString("logfile")
	if logfile == "" {
		fmt.Println("Error: `logfile` is not defined in config")
		os.Exit(2)
	}
	mkdir_p_logdir(logfile)

	// - log --------------------
	// f, err := os.OpenFile(getLogFilePath(), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	f, err := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v \n", err)
	}
	defer f.Close()
	log.SetOutput(f)
	// -------------------- log -

	srv := servicelib.NewService(config.APP_NAME, config.APP_DESC)

	if len(os.Args) >= 2 {
		cmd := strings.ToLower(os.Args[1])
		err = servicelib.HandleCmd(srv, cmd)
		if err != nil {
			usage(fmt.Sprintf("failed to %s %s: %v", cmd, config.APP_NAME, err))
		}
	} else {
		isIntSess, err := srv.IsAnInteractiveSession()
		if err != nil {
			log.Fatalf("failed to determine if we are running in an interactive session: %v", err)
		}
		if !isIntSess {
			runService(config.APP_NAME, false)
			return
		}
		// runService(config.APP_NAME, false)
	}

	return
}
