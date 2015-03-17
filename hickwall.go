package main

import (
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/codegangsta/cli"
	"github.com/oliveagle/hickwall/command"
	"github.com/oliveagle/hickwall/config"
	"github.com/oliveagle/hickwall/servicelib"
	"os"
	// "strings"
)

var err error

func main() {

	// pretty.Println(config.Conf)
	// os.Exit(1)

	defer log.Flush()

	app := cli.NewApp()
	app.Name = config.APP_NAME
	app.Usage = config.APP_DESC
	app.Version = config.VERSION

	app.Commands = []cli.Command{
		//TODO: configuration test, reload
		// {
		// 	Name:      "config",
		// 	ShortName: "",
		// 	Usage:     "config",
		// 	Subcommands: []cli.Command{
		// 		{
		// 			Name:      "test",
		// 			ShortName: "",
		// 			Usage:     "test",
		// 			Action:    command.CmdConfigTest,
		// 		},
		// 		{
		// 			Name:      "reload",
		// 			ShortName: "",
		// 			Usage:     "reload",
		// 			Action:    command.CmdConfigReload,
		// 		},
		// 	},
		// },
		{
			Name:      "service",
			ShortName: "",
			Usage:     "service",
			Subcommands: []cli.Command{
				{
					Name:      "status",
					ShortName: "",
					Usage:     "status",
					Action:    command.CmdServiceStatus,
				},
				{
					Name:   "install",
					Usage:  "install service",
					Action: command.CmdServiceInstall,
				},
				{
					Name:   "remove",
					Usage:  "remove service",
					Action: command.CmdServiceRemove,
				},
				{
					Name:   "start",
					Usage:  "start service",
					Action: command.CmdServiceStart,
				},
				{
					Name:   "stop",
					Usage:  "stop service",
					Action: command.CmdServiceStop,
				},
			},
		},
	}

	if len(os.Args) >= 2 {
		app.Run(os.Args)
	} else {
		isIntSess, err := servicelib.IsAnInteractiveSession()
		if err != nil {
			log.Error("failed to determine if we are running in an interactive session or not: %v", err)
			return
		}
		if !isIntSess {
			fmt.Println("Running ... ")
			runService(false)
			return
		}
	}

	return
}
