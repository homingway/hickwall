package command

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/oliveagle/hickwall/servicelib"
	log "github.com/oliveagle/seelog"
	"os"
	// "sync"
)

func CmdServiceStatus(c *cli.Context) {
	log.Debug("CmdServiceStatus")

	// -----------------------------------
	state, err := HelperService.Status()
	if err != nil {
		fmt.Println("error: ", err)
		log.Debug("error: ", err)
	} else {
		fmt.Printf("service %s is %s\n", HelperService.Name(), servicelib.StateToString(state))
		log.Debugf("service %s is %s\n", HelperService.Name(), servicelib.StateToString(state))
	}

	state, err = PrimaryService.Status()
	if err != nil {
		fmt.Println("error: ", err)
		log.Debug("error: ", err)
	} else {
		fmt.Printf("service %s is %s\n", PrimaryService.Name(), servicelib.StateToString(state))
		log.Debugf("service %s is %s\n", PrimaryService.Name(), servicelib.StateToString(state))
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
	log.Debug("CmdServiceInstall")
	// -----------------------------------
	err := HelperService.InstallService()
	if err != nil {
		fmt.Println("error: ", err)
		log.Debug("error: ", err)
	} else {
		fmt.Printf("service %s installed\n", HelperService.Name())
		log.Debugf("service %s installed\n", HelperService.Name())
	}

	err = PrimaryService.InstallService()
	if err != nil {
		fmt.Println("error: ", err)
		log.Debug("error: ", err)
	} else {
		fmt.Printf("service %s installed\n", PrimaryService.Name())
		log.Debugf("service %s installed\n", PrimaryService.Name())
	}

}

func CmdServiceRemove(c *cli.Context) {
	log.Debug("CmdServiceRemove")

	// -----------------------------------
	err := HelperService.RemoveService()
	if err != nil {
		fmt.Println("error: ", err)
		log.Debug("error: ", err)
	} else {
		fmt.Printf("service %s removed\n", HelperService.Name())
		log.Debugf("service %s removed\n", HelperService.Name())
	}

	err = PrimaryService.RemoveService()
	if err != nil {
		fmt.Println("error: ", err)
		log.Debug("error: ", err)
	} else {
		fmt.Printf("service %s removed\n", PrimaryService.Name())
		log.Debugf("service %s removed\n", PrimaryService.Name())
	}

}

func CmdServiceStart(c *cli.Context) {
	log.Debug("CmdServiceStart")
	// -----------------------------------
	err := HelperService.StartService()
	if err != nil {
		fmt.Println("error: ", err)
		log.Debug("error: ", err)
	} else {
		fmt.Printf("service %s started\n", HelperService.Name())
		log.Debugf("service %s started\n", HelperService.Name())
	}

	err = PrimaryService.StartService()
	if err != nil {
		fmt.Println("error: ", err)
		log.Debug("error: ", err)
	} else {
		fmt.Printf("service %s started\n", PrimaryService.Name())
		log.Debugf("service %s started\n", PrimaryService.Name())
	}
}

func CmdServiceStop(c *cli.Context) {
	log.Debug("CmdServiceStop")

	err := HelperService.StopService()
	if err != nil {
		fmt.Println("error: ", err)
		log.Debug("error: ", err)
	} else {
		fmt.Printf("service %s stopped\n", HelperService.Name())
		log.Debugf("service %s stopped\n", HelperService.Name())
	}

	err = PrimaryService.StopService()
	if err != nil {
		fmt.Println("error: ", err)
		log.Debug("error: ", err)
	} else {
		fmt.Printf("service %s stopped\n", PrimaryService.Name())
		log.Debugf("service %s stopped\n", PrimaryService.Name())
	}
}

func CmdServiceRestart(c *cli.Context) {
	log.Debug("CmdServiceRestart")

	CmdServiceStop(c)
	CmdServiceStart(c)
}
