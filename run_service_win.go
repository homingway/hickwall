// +build windows

package main

import (
	"code.google.com/p/winsvc/svc"
	// "fmt"
	// log "github.com/cihub/seelog"
	log "github.com/oliveagle/hickwall/_third_party/seelog"
	"time"
)

type myservice struct{}

func (this *myservice) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	log.Info("myservice.Execute\r\n")
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown | svc.AcceptPauseAndContinue
	changes <- svc.Status{State: svc.StartPending}

	fasttick := time.Tick(500 * time.Millisecond)
	slowtick := time.Tick(2 * time.Second)
	tick := fasttick

	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

	// major loop for signal processing.
loop:
	for {
		select {
		case <-tick:
			log.Info("beep")
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
				tick = slowtick
			case svc.Continue:
				changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
				tick = fasttick
			default:
				log.Error("unexpected control request #%d", c)
			}
		}
	}
	changes <- svc.Status{State: svc.StopPending}
	return
}

func runService(name string, isDebug bool) {
	log.Debug("runService: starting %s service \r\n", name)
	err := svc.Run(name, &myservice{})
	if err != nil {
		log.Debug("runService: Error: %s service failed: %v\r\n", name, err)
		return
	}
	log.Debug("runService: %s service stopped\r\n", name)
}
