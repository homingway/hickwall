package command

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/oliveagle/hickwall/logging"
	"github.com/oliveagle/hickwall/servicelib"
	"os"
)

func CmdServiceStatus(c *cli.Context) {
	logging.Trace("CmdServiceStatus")

	// -----------------------------------
	state, err := HelperService.Status()
	if err != nil {
		fmt.Println("error: ", err)
		logging.Tracef("error: %v", err)
	} else {
		fmt.Printf("service %s is %s\n", HelperService.Name(), servicelib.StateToString(state))
		logging.Tracef("service %s is %s\n", HelperService.Name(), servicelib.StateToString(state))
	}

	state, err = PrimaryService.Status()
	if err != nil {
		fmt.Println("error: ", err)
		logging.Tracef("error: %v", err)

	} else {
		fmt.Printf("service %s is %s\n", PrimaryService.Name(), servicelib.StateToString(state))
		logging.Tracef("service %s is %s\n", PrimaryService.Name(), servicelib.StateToString(state))
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
	logging.Trace("CmdServiceInstall")
	// -----------------------------------
	err := HelperService.InstallService()
	if err != nil {
		fmt.Println("error: ", err)
		logging.Tracef("error: %v", err)
	} else {
		fmt.Printf("service %s installed\n", HelperService.Name())
		logging.Tracef("service %s installed\n", HelperService.Name())
	}

	err = PrimaryService.InstallService()
	if err != nil {
		fmt.Println("error: ", err)
		logging.Tracef("error: %v", err)
	} else {
		fmt.Printf("service %s installed\n", PrimaryService.Name())
		logging.Tracef("service %s installed\n", PrimaryService.Name())
	}

}

func CmdServiceRemove(c *cli.Context) {
	logging.Trace("CmdServiceRemove")

	// -----------------------------------
	err := HelperService.RemoveService()
	if err != nil {
		fmt.Println("error: ", err)
		logging.Tracef("error: %v", err)
	} else {
		fmt.Printf("service %s removed\n", HelperService.Name())
		logging.Tracef("service %s removed\n", HelperService.Name())
	}

	err = PrimaryService.RemoveService()
	if err != nil {
		fmt.Println("error: ", err)
		logging.Tracef("error: %v", err)
	} else {
		fmt.Printf("service %s removed\n", PrimaryService.Name())
		logging.Tracef("service %s removed\n", PrimaryService.Name())
	}

}

func CmdServiceStart(c *cli.Context) {
	logging.Trace("CmdServiceStart")
	// -----------------------------------
	err := HelperService.StartService()
	if err != nil {
		fmt.Println("error: ", err)
		logging.Tracef("error: %v", err)
	} else {
		fmt.Printf("service %s started\n", HelperService.Name())
		logging.Tracef("service %s started\n", HelperService.Name())
	}

	err = PrimaryService.StartService()
	if err != nil {
		fmt.Println("error: ", err)
		logging.Tracef("error: %v", err)
	} else {
		fmt.Printf("service %s started\n", PrimaryService.Name())
		logging.Tracef("service %s started\n", PrimaryService.Name())
	}
}

func CmdServiceStop(c *cli.Context) {
	logging.Trace("CmdServiceStop")

	err := HelperService.StopService()
	if err != nil {
		fmt.Println("error: ", err)
		logging.Tracef("error: %v", err)
	} else {
		fmt.Printf("service %s stopped\n", HelperService.Name())
		logging.Tracef("service %s stopped\n", HelperService.Name())
	}

	err = PrimaryService.StopService()
	if err != nil {
		fmt.Println("error: ", err)
		logging.Tracef("error: %v", err)
	} else {
		fmt.Printf("service %s stopped\n", PrimaryService.Name())
		logging.Tracef("service %s stopped\n", PrimaryService.Name())
	}
}

func CmdServiceRestart(c *cli.Context) {
	logging.Trace("CmdServiceRestart")

	CmdServiceStop(c)
	CmdServiceStart(c)
}
