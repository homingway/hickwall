package command

import (
	// "fmt"
	log "github.com/cihub/seelog"
	"github.com/codegangsta/cli"
	"github.com/oliveagle/hickwall/servicelib"
	"os"
)

func CmdServiceStatus(c *cli.Context) {
	log.Debug("check service status: ", PrimaryService.Name())
	state, err := PrimaryService.Status()
	if err != nil {
		log.Error(err)
	} else {
		log.Infof("service %s is %s\n", PrimaryService.Name(), servicelib.StateToString(state))
	}
	// -----------------------------------
	log.Debugf("check service status: ", HelperService.Name())
	state, err = HelperService.Status()
	if err != nil {
		log.Error(err)
	} else {
		log.Infof("service %s is %s\n", HelperService.Name(), servicelib.StateToString(state))
	}
}

func CmdServiceStatusCode(c *cli.Context) {
	// inno setup need to know service running state.
	state, err := PrimaryService.Status()
	if err != nil {
		os.Exit(0)
	} else {
		os.Exit(int(state))
	}
}

func CmdServiceInstall(c *cli.Context) {
	log.Debug("Installing service: ", PrimaryService.Name())
	err := PrimaryService.InstallService()
	if err != nil {
		log.Error(err)
	} else {
		log.Infof("service %s installed", PrimaryService.Name())
	}
	// -----------------------------------
	log.Debug("Installing service: ", HelperService.Name())
	err = HelperService.InstallService()
	if err != nil {
		log.Error(err)
	} else {
		log.Infof("service %s installed", HelperService.Name())
	}

}

func CmdServiceRemove(c *cli.Context) {
	log.Debug("removing service: ", PrimaryService.Name())
	err := PrimaryService.RemoveService()
	if err != nil {
		log.Error(err)
	} else {
		log.Infof("service %s removed", PrimaryService.Name())
	}
	// -----------------------------------
	log.Debug("removing service: ", HelperService.Name())
	err = HelperService.RemoveService()
	if err != nil {
		log.Error(err)
	} else {
		log.Infof("service %s removed", HelperService.Name())
	}

}

func CmdServiceStart(c *cli.Context) {
	log.Debug("starting service: ", PrimaryService.Name())
	err := PrimaryService.StartService()
	if err != nil {
		log.Error(err)
	} else {
		log.Infof("service %s started", PrimaryService.Name())
	}
	// -----------------------------------
	log.Debug("starting service: ", HelperService.Name())
	err = HelperService.StartService()
	if err != nil {
		log.Error(err)
	} else {
		log.Infof("service %s started", HelperService.Name())
	}
}

func CmdServiceStop(c *cli.Context) {
	log.Debug("stopping service: ", PrimaryService.Name())
	err := PrimaryService.StopService()
	if err != nil {
		log.Error(err)
	} else {
		log.Infof("service %s stopped", PrimaryService.Name())
	}

	// -----------------------------------
	log.Debug("stopping service: ", HelperService.Name())
	err = HelperService.StopService()
	if err != nil {
		log.Error(err)
	} else {
		log.Infof("service %s stopped", HelperService.Name())
	}
}

func CmdServicePause(c *cli.Context) {
	log.Debug("pausing service: ", PrimaryService.Name())
	err := PrimaryService.PauseService()
	if err != nil {
		log.Error(err)
	} else {
		log.Infof("service %s paused", PrimaryService.Name())
	}
}

func CmdServiceContinue(c *cli.Context) {
	log.Debug("continuing service: ", PrimaryService.Name())
	err := PrimaryService.ContinueService()
	if err != nil {
		log.Error(err)
	} else {
		log.Infof("service %s continued", PrimaryService.Name())
	}
}
