// +build linux darwin

package servicelib

import (
	"fmt"
	// "github.com/oliveagle/hickwall/config"
	log "github.com/oliveagle/seelog"
	// "github.com/op/go-logging"
	// "github.com/VividCortex/robustly"
)

func printCmdRes(str string, err error) {
	fmt.Println(str)
}

func IsAnInteractiveSession() (bool, error) {
	return false, nil
}

func (this *Service) InstallService() error {
	log.Debug("ServiceManager.InstallService")

	str, err := this.Install()
	printCmdRes(str, err)
	return err
}

func (this *Service) RemoveService() error {
	log.Debug("ServiceManager.RemoveService")

	str, err := this.Remove()
	printCmdRes(str, err)
	return err
}

func (this *Service) Status() (State, error) {
	log.Error("ServiceManagement.Status not supported")
	return Unknown, fmt.Errorf("ServerMangement.Status not supported")
}

func (this *Service) StartService() error {
	log.Debug("ServiceManager.StartService")

	str, err := this.Start()
	printCmdRes(str, err)
	return err
}

func (this *Service) StopService() error {
	log.Debug("ServiceManager.StopService")

	str, err := this.Stop()
	printCmdRes(str, err)
	return err
}

func (this *Service) PauseService() error {
	log.Error("ServiceManager.PauseServicen not supported ")

	return nil
}

func (this *Service) ContinueService() error {
	log.Error("ServiceManager.ContinueService not supported ")
	return nil
}
