package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/oliveagle/hickwall/command"
	"github.com/oliveagle/hickwall/logging"
	"github.com/oliveagle/hickwall/servicelib"
	"github.com/oliveagle/hickwall/utils"
	"os"
)

var err error
var Build string
var Version string

func main() {
	defer utils.Recover_and_log()
	logging.Debug("hickwall main ------- sub packages init process finished")

	app := cli.NewApp()
	app.Name = "hickwall"
	app.Usage = "monitoring system"
	app.Version = fmt.Sprintf("%s - %s", Version, Build)

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
			ShortName: "s",
			Usage:     "service",
			Subcommands: []cli.Command{
				{
					Name:      "status",
					ShortName: "s",
					Usage:     "status",
					Action:    command.CmdServiceStatus,
				},
				{
					Name:   "statuscode",
					Usage:  "statuscode(internal use only.)",
					Action: command.CmdServiceStatusCode,
				},
				{
					Name:      "install",
					ShortName: "i",
					Usage:     "install service",
					Action:    command.CmdServiceInstall,
				},
				{
					Name:      "remove",
					ShortName: "d",
					Usage:     "remove service",
					Action:    command.CmdServiceRemove,
				},
				{
					Name:      "start",
					ShortName: "g",
					Usage:     "start service.",
					Action:    command.CmdServiceStart,
				},
				{
					Name:      "stop",
					ShortName: "x",
					Usage:     "stop service.",
					Action:    command.CmdServiceStop,
				},
				{
					Name:      "restart",
					ShortName: "n",
					Usage:     "restart service",
					Action:    command.CmdServiceRestart,
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
				run(false, false)
			},
		},
		{
			Name:   "config",
			Usage:  "show config info",
			Action: command.CmdShowConfig,
		},
	}

	if len(os.Args) >= 2 {
		app.Run(os.Args)
	} else {
		logging.Debug("len os.args < 2")

		isIntSess, err := servicelib.IsAnInteractiveSession()
		if err != nil {
			logging.Errorf("failed to determine if we are running in an interactive session or not: %v", err)
			return
		}

		if !isIntSess {
			logging.Debug("running as service ... ")
			run(false, true)
			return
		}

		//print help here.
		app.Run(os.Args)
	}
	return
}
