package command

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/oliveagle/hickwall/servicelib"
	"os"
	// "sync"
	"github.com/oliveagle/hickwall/logging"
)

func CmdServiceStatus(c *cli.Context) {
	logging.Infof("CmdServiceStatus")

	// -----------------------------------
	state, err := HelperService.Status()
	if err != nil {
		fmt.Println("error: ", err)
		logging.Errorf("error: %v", err)
	} else {
		fmt.Printf("service %s is %s\n", HelperService.Name(), servicelib.StateToString(state))
		logging.Debugf("service %s is %s\n", HelperService.Name(), servicelib.StateToString(state))
	}

	state, err = PrimaryService.Status()
	if err != nil {
		fmt.Println("error: ", err)
		logging.Errorf("error: %v", err)

	} else {
		fmt.Printf("service %s is %s\n", PrimaryService.Name(), servicelib.StateToString(state))
		logging.Debugf("service %s is %s\n", PrimaryService.Name(), servicelib.StateToString(state))
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
	logging.Debug("CmdServiceInstall")
	// -----------------------------------
	err := HelperService.InstallService()
	if err != nil {
		fmt.Println("error: ", err)
		logging.Debugf("error: %v", err)
	} else {
		fmt.Printf("service %s installed\n", HelperService.Name())
		logging.Debugf("service %s installed\n", HelperService.Name())
	}

	err = PrimaryService.InstallService()
	if err != nil {
		fmt.Println("error: ", err)
		logging.Debugf("error: %v", err)
	} else {
		fmt.Printf("service %s installed\n", PrimaryService.Name())
		logging.Debugf("service %s installed\n", PrimaryService.Name())
	}

}

func CmdServiceRemove(c *cli.Context) {
	logging.Debug("CmdServiceRemove")

	// -----------------------------------
	err := HelperService.RemoveService()
	if err != nil {
		fmt.Println("error: ", err)
		logging.Debugf("error: %v", err)
	} else {
		fmt.Printf("service %s removed\n", HelperService.Name())
		logging.Debugf("service %s removed\n", HelperService.Name())
	}

	err = PrimaryService.RemoveService()
	if err != nil {
		fmt.Println("error: ", err)
		logging.Debugf("error: %v", err)
	} else {
		fmt.Printf("service %s removed\n", PrimaryService.Name())
		logging.Debugf("service %s removed\n", PrimaryService.Name())
	}

}

func CmdServiceStart(c *cli.Context) {
	logging.Debug("CmdServiceStart")
	// -----------------------------------
	err := HelperService.StartService()
	if err != nil {
		fmt.Println("error: ", err)
		logging.Debugf("error: %v", err)
	} else {
		fmt.Printf("service %s started\n", HelperService.Name())
		logging.Debugf("service %s started\n", HelperService.Name())
	}

	err = PrimaryService.StartService()
	if err != nil {
		fmt.Println("error: ", err)
		logging.Debugf("error: %v", err)
	} else {
		fmt.Printf("service %s started\n", PrimaryService.Name())
		logging.Debugf("service %s started\n", PrimaryService.Name())
	}
}

func CmdServiceStop(c *cli.Context) {
	logging.Debug("CmdServiceStop")

	err := HelperService.StopService()
	if err != nil {
		fmt.Println("error: ", err)
		logging.Debugf("error: %v", err)
	} else {
		fmt.Printf("service %s stopped\n", HelperService.Name())
		logging.Debugf("service %s stopped\n", HelperService.Name())
	}

	err = PrimaryService.StopService()
	if err != nil {
		fmt.Println("error: ", err)
		logging.Debugf("error: %v", err)
	} else {
		fmt.Printf("service %s stopped\n", PrimaryService.Name())
		logging.Debugf("service %s stopped\n", PrimaryService.Name())
	}
}

func CmdServiceRestart(c *cli.Context) {
	logging.Debug("CmdServiceRestart")

	CmdServiceStop(c)
	CmdServiceStart(c)
}
