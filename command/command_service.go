package command

import (
	// "fmt"
	"github.com/codegangsta/cli"
	log "github.com/oliveagle/hickwall/_third_party/seelog"
	"github.com/oliveagle/hickwall/config"
	"github.com/oliveagle/hickwall/servicelib"
)

var service = servicelib.NewService(config.APP_NAME, config.APP_DESC)

func CmdServiceStatus(c *cli.Context) {
	log.Info("CmdServiceStatus")
	err := service.Status()
	if err != nil {
		log.Errorf("CmdServiceStatus: %v", err)
	}
}

func CmdServiceInstall(c *cli.Context) {
	log.Info("CmdServiceInstall")
	err := service.InstallService()
	if err != nil {
		log.Errorf("CmdServiceStatus: %v", err)
	}
}

func CmdServiceRemove(c *cli.Context) {
	log.Info("CmdServiceRemove")
	err := service.RemoveService()
	if err != nil {
		log.Errorf("CmdServiceStatus: %v", err)
	}
}

func CmdServiceStart(c *cli.Context) {
	log.Info("CmdServiceStart")
	err := service.StartService()
	if err != nil {
		log.Errorf("CmdServiceStatus: %v", err)
	}
}

func CmdServiceStop(c *cli.Context) {
	log.Info("CmdServiceStop")
	err := service.StopService()
	if err != nil {
		log.Errorf("CmdServiceStatus: %v", err)
	}
}

func CmdServicePause(c *cli.Context) {
	log.Info("CmdServicePause")
	err := service.PauseService()
	if err != nil {
		log.Errorf("CmdServiceStatus: %v", err)
	}
}

func CmdServiceContinue(c *cli.Context) {
	log.Info("CmdServiceContinue")
	err := service.ContinueService()
	if err != nil {
		log.Errorf("CmdServiceStatus: %v", err)
	}
}
