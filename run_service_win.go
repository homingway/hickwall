// +build windows

package main

import (
	"code.google.com/p/winsvc/svc"
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/oliveagle/go-collectors/datapoint"
	"github.com/oliveagle/hickwall/backends"
	"github.com/oliveagle/hickwall/collectors"
	"time"

	"github.com/oliveagle/hickwall/utils"
)

type myservice struct{}

func (this *myservice) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	log.Info("myservice.Execute\r\n")
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown | svc.AcceptPauseAndContinue
	changes <- svc.Status{State: svc.StartPending}

	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

	mdCh := make(chan *datapoint.MultiDataPoint)

	collectors.RunAllCollectors(mdCh)
	backends.RunBackends()
	defer backends.CloseBackends()

	utils.HttpPprofServe(6060)

	tick := time.Tick(time.Second * time.Duration(1))
	go func() {
		for {
			select {
			case <-tick:
				md, _ := collectors.C_hickwall(nil)
				mdCh <- md
			}
		}
	}()

	// major loop for signal processing.
loop:
	for {
		select {
		case md, err := <-mdCh:
			fmt.Println("MultiDataPoint: ", md, err)
			for _, p := range *md {
				fmt.Println(" point ---> ", p)
			}
			backends.WriteToBackends(*md)
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
				// testing deadlock from https://code.google.com/p/winsvc/issues/detail?id=4
				time.Sleep(100 * time.Millisecond)
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				break loop
			case svc.Pause:
				changes <- svc.Status{State: svc.Paused, Accepts: cmdsAccepted}
				log.Error("win.Pause not implemented yet")
			case svc.Continue:
				changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
				log.Error("win.Continue not implemented yet")
			default:
				log.Error("unexpected control request #%d", c)
			}
		}
	}
	changes <- svc.Status{State: svc.StopPending}
	return
}

func runService(name string, isDebug bool) {
	fmt.Println("hahah")
	log.Debug("runService: starting %s service \r\n", name)
	err := svc.Run(name, &myservice{})
	if err != nil {
		log.Debug("runService: Error: %s service failed: %v\r\n", name, err)
		return
	}
	log.Debug("runService: %s service stopped\r\n", name)
}
