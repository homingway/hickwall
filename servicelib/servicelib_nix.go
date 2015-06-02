// +build linux darwin

package servicelib

import (
	"fmt"
	"github.com/oliveagle/hickwall/logging"
)

func printCmdRes(str string, err error) {
	fmt.Println(str)
}

func IsAnInteractiveSession() (bool, error) {
	return false, nil
}

func (this *Service) InstallService() error {
	logging.Debug("ServiceManager.InstallService")

	str, err := this.Install()
	printCmdRes(str, err)
	return err
}

func (this *Service) RemoveService() error {
	logging.Debug("ServiceManager.RemoveService")

	str, err := this.Remove()
	printCmdRes(str, err)
	return err
}

func (this *Service) Status() (State, error) {
	logging.Error("ServiceManagement.Status not supported")
	return Unknown, fmt.Errorf("ServerMangement.Status not supported")
}

func (this *Service) StartService() error {
	logging.Debug("ServiceManager.StartService")

	str, err := this.Start()
	printCmdRes(str, err)
	return err
}

func (this *Service) StopService() error {
	logging.Debug("ServiceManager.StopService")

	str, err := this.Stop()
	printCmdRes(str, err)
	return err
}

func (this *Service) PauseService() error {
	logging.Error("ServiceManager.PauseServicen not supported ")

	return nil
}

func (this *Service) ContinueService() error {
	logging.Error("ServiceManager.ContinueService not supported ")
	return nil
}
