// +build windows

package main

import (
	"code.google.com/p/winsvc/eventlog"
	"code.google.com/p/winsvc/svc"
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/oliveagle/go-collectors/datapoint"
	"github.com/oliveagle/hickwall/backends"
	"github.com/oliveagle/hickwall/collectors"
	"github.com/oliveagle/hickwall/command"
	"github.com/oliveagle/hickwall/config"
	"github.com/oliveagle/hickwall/servicelib"
	"github.com/oliveagle/hickwall/utils"
	"time"
)

func start_service_if_stopped(service *servicelib.Service) {
	state, err := service.Status()
	if err != nil {
		log.Errorf("CmdServiceStatus: %v\n", err)
		return
	}
	if state == servicelib.Stopped {
		log.Infof("service %s is stopped! trying to start service again", service.Name())

		err := service.StartService()
		if err != nil {
			log.Info("start service failed: ", err)
		} else {
			log.Info("service started. ")
		}
	} else {
		log.Info("Serivce state: ", servicelib.StateToString(state))
	}
}

type serviceHandler struct{}

func runAsPrimaryService(elog *eventlog.Log, args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	elog.Info(1, "runAsPrimaryService -- 1 --")
	log.Info("runAsPrimaryService -- 1 --")

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
				go start_service_if_stopped(command.HelperService)
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
				log.Info("svc.Stop or svc.Shutdown is triggered")
				break loop
			case svc.Pause:
				changes <- svc.Status{State: svc.Paused, Accepts: cmdsAccepted}
				log.Info("win.Pause not implemented yet")
			case svc.Continue:
				changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
				log.Info("win.Continue not implemented yet")
			default:
				log.Info("unexpected control request #%d", c)
			}
		}
	}
	changes <- svc.Status{State: svc.StopPending}
	log.Info("serviceHandler stopped")
	return
}

func runAsHelperService(elog *eventlog.Log, args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	log.Info("runAsHelperService -- 2 --")

	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown
	changes <- svc.Status{State: svc.StartPending}

	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

	tick := time.Tick(time.Second * time.Duration(1))

	// major loop for signal processing.
loop:
	for {
		select {
		case <-tick:
			go start_service_if_stopped(command.PrimaryService)
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
				// testing deadlock from https://code.google.com/p/winsvc/issues/detail?id=4
				time.Sleep(100 * time.Millisecond)
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				log.Info("svc.Stop or svc.Shutdown is triggered")
				break loop
			default:
				log.Info("unexpected control request #%d", c)
			}
		}
	}
	changes <- svc.Status{State: svc.StopPending}
	log.Info("help_service stopped")
	return

}

func (this *serviceHandler) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	elog, err := eventlog.Open("hickwall")
	if err != nil {
		log.Error("Cannot open eventlog: hickwall")
		return
	}
	defer elog.Close()
	elog.Info(1, "serviceHandler.Execute")

	err = config.Init()
	if err != nil {
		elog.Error(3, fmt.Sprintf("config.Init Failed: %v", err))
		return
	}
	elog.Info(1, "Config.Init")

	log.Error("serviceHandler.Execute: args:", args)

	if len(args) > 0 {
		svc_name := args[0]
		if svc_name == "hickwall" {
			return runAsPrimaryService(elog, args, r, changes)
		} else {
			return runAsHelperService(elog, args, r, changes)
		}
	}

	return runAsPrimaryService(elog, args, r, changes)
}
func runService(isDebug bool) {
	err := svc.Run(command.PrimaryService.Name(), &serviceHandler{})
	if err != nil {
		log.Error("runService: failed: %v\r\n", err)
	}
}
