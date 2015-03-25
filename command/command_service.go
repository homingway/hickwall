package command

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/oliveagle/hickwall/servicelib"
	"os"
	"sync"
)

func CmdServiceStatus(c *cli.Context) {
	state, err := PrimaryService.Status()
	if err != nil {
		fmt.Println("error: ", err)
	} else {
		fmt.Printf("service %s is %s\n", PrimaryService.Name(), servicelib.StateToString(state))
	}
	// -----------------------------------
	state, err = HelperService.Status()
	if err != nil {
		fmt.Println("error: ", err)
	} else {
		fmt.Printf("service %s is %s\n", HelperService.Name(), servicelib.StateToString(state))
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
	err := PrimaryService.InstallService()
	if err != nil {
		fmt.Println("error: ", err)
	} else {
		fmt.Printf("service %s installed\n", PrimaryService.Name())
	}
	// -----------------------------------
	err = HelperService.InstallService()
	if err != nil {
		fmt.Println("error: ", err)
	} else {
		fmt.Printf("service %s installed\n", HelperService.Name())
	}

}

func CmdServiceRemove(c *cli.Context) {
	err := PrimaryService.RemoveService()
	if err != nil {
		fmt.Println("error: ", err)
	} else {
		fmt.Printf("service %s removed\n", PrimaryService.Name())
	}
	// -----------------------------------
	err = HelperService.RemoveService()
	if err != nil {
		fmt.Println("error: ", err)
	} else {
		fmt.Printf("service %s removed\n", HelperService.Name())
	}

}

func CmdServiceStart(c *cli.Context) {
	err := PrimaryService.StartService()
	if err != nil {
		fmt.Println("error: ", err)
	} else {
		fmt.Printf("service %s started\n", PrimaryService.Name())
	}
	// -----------------------------------
	err = HelperService.StartService()
	if err != nil {
		fmt.Println("error: ", err)
	} else {
		fmt.Printf("service %s started\n", HelperService.Name())
	}
}

func CmdServiceStop(c *cli.Context) {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		err := PrimaryService.StopService()
		if err != nil {
			fmt.Println("error: ", err)
		} else {
			fmt.Printf("service %s stopped\n", PrimaryService.Name())
		}
	}()

	// -----------------------------------
	wg.Add(1)
	go func() {
		defer wg.Done()

		err := HelperService.StopService()
		if err != nil {
			fmt.Println("error: ", err)
		} else {
			fmt.Printf("service %s stopped\n", HelperService.Name())
		}
	}()

	wg.Wait()
}

func CmdServiceRestart(c *cli.Context) {
	CmdServiceStop(c)
	CmdServiceStart(c)
}

func CmdServicePause(c *cli.Context) {
	err := PrimaryService.PauseService()
	if err != nil {
		fmt.Println("error: ", err)
	} else {
		fmt.Printf("service %s paused\n", PrimaryService.Name())
	}
}

func CmdServiceContinue(c *cli.Context) {
	err := PrimaryService.ContinueService()
	if err != nil {
		fmt.Println("error: ", err)
	} else {
		fmt.Printf("service %s continued\n", PrimaryService.Name())
	}
}
