// +build windows

package main

import (
	"code.google.com/p/winsvc/eventlog"
	"code.google.com/p/winsvc/svc"
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/oliveagle/hickwall/backends"
	"github.com/oliveagle/hickwall/collectorlib"
	"github.com/oliveagle/hickwall/collectors"
	"github.com/oliveagle/hickwall/command"
	"github.com/oliveagle/hickwall/config"
	"github.com/oliveagle/hickwall/servicelib"
	"github.com/oliveagle/hickwall/utils"
	"time"
)

func start_service_if_stopped(elog *eventlog.Log, service *servicelib.Service) {
	state, err := service.Status()
	if err != nil {
		elog.Error(3, fmt.Sprintf("CmdServiceStatus Error: %v", err))
		return
	}
	if state == servicelib.Stopped {
		elog.Warning(2, fmt.Sprintf("service %s is stopped! trying to start service again", service.Name()))

		err := service.StartService()
		if err != nil {
			elog.Error(3, fmt.Sprintf("start service failed: ", err))
		} else {
			elog.Info(1, fmt.Sprintf("service %s started", service.Name()))
		}
	}
}

type serviceHandler struct{}

func log_rush() {
	log.Info(`A Dream Within A Dream - by by Edgar Allan Poe
		Take this kiss upon the brow!
		And, in parting from you now,
		Thus much let me avow--
		You are not wrong, who deem
		That my days have been a dream;
		Yet if hope has flown away
		In a night, or in a day,
		In a vision, or in none,
		Is it therefore the less gone?
		All that we see or seem
		Is but a dream within a dream.
		I stand amid the roar
		Of a surf-tormented shore,
		And I hold within my hand
		Grains of the golden sand--
		How few! yet how they creep
		Through my fingers to the deep,
		While I weep--while I weep!
		O God! can I not grasp
		Them with a tighter clasp?
		O God! can I not save
		One from the pitiless wave?
		Is all that we see or seem
		But a dream within a dream?`)
	log.Info(`Stopping by Woods on a Snowy Evening - by by Robert Frost
		Whose woods these are I think I know.
		His house is in the village, though;
		He will not see me stopping here
		To watch his woods fill up with snow.
		My little horse must think it queer
		To stop without a farmhouse near
		Between the woods and frozen lake
		The darkest evening of the year.

		He gives his harness bells a shake
		To ask if there is some mistake.
		The only other sound's the sweep
		Of easy wind and downy flake.
		The woods are lovely, dark and deep,
		But I have promises to keep,
		And miles to go before I sleep,
		And miles to go before I sleep.`)
	log.Info(`A Dream Within A Dream - by by Edgar Allan Poe
		Take this kiss upon the brow!
		And, in parting from you now,
		Thus much let me avow--
		You are not wrong, who deem
		That my days have been a dream;
		Yet if hope has flown away
		In a night, or in a day,
		In a vision, or in none,
		Is it therefore the less gone?
		All that we see or seem
		Is but a dream within a dream.
		I stand amid the roar
		Of a surf-tormented shore,
		And I hold within my hand
		Grains of the golden sand--
		How few! yet how they creep
		Through my fingers to the deep,
		While I weep--while I weep!
		O God! can I not grasp
		Them with a tighter clasp?
		O God! can I not save
		One from the pitiless wave?
		Is all that we see or seem
		But a dream within a dream?`)
}

func runAsPrimaryService(elog *eventlog.Log, args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {

	err = config.Init()
	if err != nil {
		elog.Error(3, fmt.Sprintf("config.Init Failed: %v", err))
		return
	}

	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown | svc.AcceptPauseAndContinue
	changes <- svc.Status{State: svc.StartPending}

	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

	mdCh := make(chan collectorlib.MultiDataPoint)

	collectors.RunAllCollectors(mdCh)
	backends.RunBackends()
	defer backends.CloseBackends()

	utils.HttpPprofServe(6060)

	tick := time.Tick(time.Second * time.Duration(1))
	go func() {
		for {
			select {
			case <-tick:
				go start_service_if_stopped(elog, command.HelperService)
				md, _ := collectors.C_hickwall(nil)
				mdCh <- md
				// log_rush()
				// log.Info("hahahah running")
			}
		}
	}()

	// major loop for signal processing.
loop:
	for {
		select {
		case md, _ := <-mdCh:
			for _, p := range md {
				log.Debug(" point ---> ", p)
			}
			backends.WriteToBackends(md)
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
				// testing deadlock from https://code.google.com/p/winsvc/issues/detail?id=4
				time.Sleep(100 * time.Millisecond)
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				elog.Info(1, "svc.Stop or svc.Shutdown is triggered")
				break loop
			case svc.Pause:
				changes <- svc.Status{State: svc.Paused, Accepts: cmdsAccepted}
				elog.Info(1, "svc.Pause not implemented yet")
			case svc.Continue:
				changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
				elog.Info(1, "svc.Continue not implemented yet")
			default:
				elog.Error(3, fmt.Sprintf("unexpected control request #%d", c))
			}
		}
	}
	changes <- svc.Status{State: svc.StopPending}
	elog.Info(1, "serviceHandler stopped")
	return
}

func runAsHelperService(elog *eventlog.Log, args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	// NOTE: helper service should not write log to file. otherwise, multiple process write to same log file will cause log
	// rotate have unexpected behaviors.

	elog.Info(1, "helper service started")

	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown
	changes <- svc.Status{State: svc.StartPending}

	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

	tick := time.Tick(time.Second * time.Duration(1))

	// major loop for signal processing.
loop:
	for {
		select {
		case <-tick:
			go start_service_if_stopped(elog, command.PrimaryService)
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
				// testing deadlock from https://code.google.com/p/winsvc/issues/detail?id=4
				time.Sleep(100 * time.Millisecond)
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				elog.Info(1, "svc.Stop or svc.Shutdown is triggered")
				break loop
			default:
				elog.Error(3, fmt.Sprintf("unexpected control request #%d", c))

			}
		}
	}
	changes <- svc.Status{State: svc.StopPending}
	elog.Info(1, "helper service stopped")
	return

}

func (this *serviceHandler) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	elog, err := eventlog.Open("hickwall")
	if err != nil {
		log.Error("Cannot open eventlog: hickwall")
		return
	}
	defer elog.Close()

	elog.Info(1, fmt.Sprintf("serviceHandler.Execute: %v", args))

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
	elog, err := eventlog.Open("hickwall")
	if err == nil {
		defer elog.Close()
	}
	elog.Info(1, "runService is called")

	err = svc.Run(command.PrimaryService.Name(), &serviceHandler{})
	if err != nil && elog != nil {
		elog.Error(3, fmt.Sprintf("runService: failed: %v\r\n", err))
	}
}
