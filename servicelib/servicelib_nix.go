// +build linux darwin

package servicelib

import (
	"fmt"
	// "github.com/oliveagle/hickwall/config"
	// log "github.com/cihub/seelog"
	log "github.com/cihub/seelog"
	// "github.com/op/go-logging"
	// "github.com/VividCortex/robustly"
)

func printCmdRes(str string, err error) {
	fmt.Println(str)
	// if err != nil {
	// 	fmt.Println(err)
	// }
}

func (this *Service) IsAnInteractiveSession() (bool, error) {
	log.Debug("IsAnInteractiveSessioin")
	return false, nil
}

func (this *Service) InstallService() error {
	log.Debug("ServiceManager.InstallService")

	str, err := this.Install()
	// log.Debug("InstallService: %s, err: %v\n", str, err)
	printCmdRes(str, err)
	return err
}

func (this *Service) RemoveService() error {
	log.Debug("ServiceManager.RemoveService")

	str, err := this.Remove()
	// log.Debug("RemoveService: %s, err: %v", str, err)
	printCmdRes(str, err)
	return err
}

func (this *Service) Status() error {
	log.Error("ServiceManagement.Status not supported")

	return nil
}

func (this *Service) StartService() error {
	log.Debug("ServiceManager.StartService")

	str, err := this.Start()
	// log.Debug("StartService: %s, err: %v", str, err)
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
