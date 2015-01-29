// +build linux darwin

package servicelib

import (
	"fmt"
	// "github.com/oliveagle/hickwall/config"
	// "log"
	log "github.com/cihub/seelog"
	// "github.com/op/go-logging"
	// "github.com/VividCortex/robustly"
)

func printCmdRes(str string, err error) {
	fmt.Println("** print cmd ****", str)
	if err != nil {
		fmt.Println(err)
	}
}

func (this *Service) IsAnInteractiveSession() (bool, error) {
	log.Debug("IsAnInteractiveSessioin")
	return false, nil
}

func (this *Service) InstallService() error {
	log.Debug("ServiceManager.InstallService")
	str, err := this.Install()
	log.Debug("InstallService: %s, err: %v\n", str, err)
	printCmdRes(str, err)
	return err
}

func (this *Service) RemoveService() error {
	log.Debug("ServiceManager.RemoveService")
	str, err := this.Remove()
	log.Debug("RemoveService: %s, err: %v", str, err)
	printCmdRes(str, err)
	return err
}

func (this *Service) Status() error {
	log.Debug("ServiceManagement.Status not supported")
	log.Info("hahahh")

	log.Debug("Status ---------- @#@#@#  logging")
	log.Error("Status ---------- @#@#@#  logging 1212")
	return nil
}

func (this *Service) StartService() error {
	log.Debug("ServiceManager.StartService")

	str, err := this.Start()
	log.Debug("StartService: %s, err: %v", str, err)
	printCmdRes(str, err)
	return err
}

func (this *Service) StopService() error {
	// log.Println("ServiceManager.StopService")
	log.Debug("ServiceManager.StopService")
	// robustly.Run(func() {

	// }, nil)
	// printCmdRes(str, err)
	str, err := this.Stop()
	printCmdRes(str, err)

	log.Trace("ServiceManager.StopService ---------- @#@#@#  logging")
	log.Debug("ServiceManager.StopService ---------- @#@#@#  logging")
	log.Info("ServiceManager.StopService ---------- @#@#@#  logging")
	log.Warn("ServiceManager.StopService ---------- @#@#@#  logging")
	log.Error("ServiceManager.StopService ---------- @#@#@#  logging 1212")
	log.Critical("ServiceManager.StopService ---------- @#@#@#  logging")
	return err
}

func (this *Service) PauseService() error {
	log.Debug("ServiceManager.PauseServicen not supported ")
	return nil
}

func (this *Service) ContinueService() error {
	log.Debug("ServiceManager.ContinueService not supported ")
	return nil
}
