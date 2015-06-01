package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/oliveagle/hickwall/command"
	"github.com/oliveagle/hickwall/servicelib"
	"github.com/oliveagle/hickwall/utils"
	log "github.com/oliveagle/seelog"
	"os"
)

var err error

func main() {
	defer utils.Recover_and_log()
	log.Debug("hickwall main ---------------------------")

	app := cli.NewApp()
	app.Name = "hickwall"
	app.Usage = "monitoring system"
	app.Version = "v0.1.2"

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
					Name:      "statuscode",
					ShortName: "",
					Usage:     "statuscode",
					Action:    command.CmdServiceStatusCode,
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
				{
					Name:   "restart",
					Usage:  "restart service",
					Action: command.CmdServiceRestart,
				},
			},
		},
		{
			Name:      "version",
			ShortName: "v",
			Usage:     "show version info",
			Action: func(c *cli.Context) {
				fmt.Printf("%s version: %s\n", app.Name, app.Version)
			},
		},
		{
			Name:      "daemon",
			ShortName: "d",
			Usage:     "run as daemon",
			Action: func(c *cli.Context) {
				fmt.Println("Running as Daemon")
				runWithoutService()
			},
		},
	}

	if len(os.Args) >= 2 {
		log.Debug("len os.args >= 2")
		app.Run(os.Args)
	} else {
		log.Debug("len os.args < 2")

		isIntSess, err := servicelib.IsAnInteractiveSession()
		if err != nil {
			log.Error("failed to determine if we are running in an interactive session or not: %v", err)
			return
		}

		if !isIntSess {
			log.Debug("running as service ... ")
			runService(false)
			return
		}
	}
	return
}
