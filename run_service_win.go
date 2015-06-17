// +build windows

package main

import (
	"code.google.com/p/winsvc/svc"
	"github.com/oliveagle/hickwall/command"
	"github.com/oliveagle/hickwall/config"
	"github.com/oliveagle/hickwall/hickwall"
	"github.com/oliveagle/hickwall/logging"
	//	"github.com/oliveagle/hickwall/servicelib"
	"github.com/oliveagle/get-rss"
	"github.com/oliveagle/hickwall/utils"
	// "runtime/debug"
	"os"
	// "runtime/pprof"
	//	"github.com/davecheney/profile"
	"time"
)

var pid int
var rss_up_limit_mb = 50.0 // MB

func init() {
	pid = os.Getpid()
}

type serviceHandler struct{}

func run(isDebug bool, daemon bool) {
	if !config.IsCoreConfigLoaded() {
		err := config.LoadCoreConfig()
		if err != nil {
			logging.Critical("Failed to load core config: ", err)
			return
		}
	}

	go func() {
	loop:
		rss_mb := float64(gs.GetCurrentRSS()) / 1024 / 1024 // Mb
		logging.Tracef("current rss: %f", rss_mb)
		if rss_mb > rss_up_limit_mb {
			logging.Criticalf("Suicide. CurrentRSS above limit: %f >= %f Mb", rss_mb, rss_up_limit_mb)
			os.Exit(1)
		}
		next := time.After(time.Second)
		<-next
		goto loop
	}()

	if daemon {
		runService(isDebug)
	} else {
		runWithoutService()
	}
}

func runWithoutService() {
	var args = []string{}
	var r = make(chan svc.ChangeRequest)
	var changes = make(chan svc.Status, 0)

	go func() {
		for {
			select {
			case <-changes:
			}
		}
	}()

	runAsPrimaryService(args, r, changes)
}

func runService(isDebug bool) {
	defer utils.Recover_and_log()
	logging.Debug("runService")

	err = svc.Run(command.PrimaryService.Name(), &serviceHandler{})
	if err != nil {
		logging.Errorf("runService: failed: %v\r\n", err)
	}
}

func runAsPrimaryService(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	logging.Info("runAsPrimaryService")
	defer utils.Recover_and_log()

	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown | svc.AcceptPauseAndContinue
	changes <- svc.Status{State: svc.StartPending}
	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

	//http://localhost:6060/debug/pprof/
	// utils.HttpPprofServe(6060)

	//	after := time.After(time.Duration(8) * time.Minute)
	// f, _ := os.Create("d:\\cpu-" + strconv.Itoa(pid) + ".pprof")
	// pprof.StartCPUProfile(f)
	// defer pprof.StopCPUProfile()

	//	cfg := profile.Config{
	//		MemProfile:     true,
	//		ProfilePath:    "./pprofs/", // store profiles in current directory
	//		NoShutdownHook: true,        // do not hook SIGINT
	//	}
	//	p := profile.Start(&cfg)
	//
	//	defer p.Stop()

	// utils.StartCPUProfile()
	// defer utils.StopCPUProfile()

	// go func() {
	// 	for {
	// 		<-time.After(time.Second * time.Duration(15))
	// 		debug.FreeOSMemory()
	// 	}
	// }()

	err := hickwall.Start()
	if err != nil {
		logging.Critical("Failed To Start hickwall: %v", err)
		return
	} else {
		defer hickwall.Stop()
	}

	logging.Debug("service event handling loop started ")
	// major loop for signal processing.
loop:
	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
				// testing deadlock from https://code.google.com/p/winsvc/issues/detail?id=4
				time.Sleep(100 * time.Millisecond)
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				logging.Info("svc.Stop or svc.Shutdown is triggered")
				break loop
			case svc.Pause:
				changes <- svc.Status{State: svc.Paused, Accepts: cmdsAccepted}
				logging.Info("svc.Pause not implemented yet")
			case svc.Continue:
				changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
				logging.Info("svc.Continue not implemented yet")
			default:
				logging.Errorf("unexpected control request #%d", c)
			}
		}
	}
	changes <- svc.Status{State: svc.StopPending}
	logging.Info("serviceHandler stopped")
	return
}

func (this *serviceHandler) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	return runAsPrimaryService(args, r, changes)
}
