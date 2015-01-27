// +build linux darwin

package servicelib

import (
	// "github.com/oliveagle/hickwall/config"
	"fmt"
	"github.com/spf13/viper"
	"log"
)

func printCmdRes(str string, err error) {
	fmt.Println(str)
	if err != nil {
		fmt.Println(err)
	}
}

func (this *Service) IsAnInteractiveSession() (bool, error) {
	log.Println("IsAnInteractiveSessioin")
	return false, nil
}

func (this *Service) InstallService() error {
	log.Println("ServiceManager.InstallService")
	str, err := this.Install()
	log.Printf("InstallService: %s, err: %v\n", str, err)
	printCmdRes(str, err)
	return err
}

func (this *Service) RemoveService() error {
	log.Println("ServiceManager.RemoveService")
	str, err := this.Remove()
	log.Println("RemoveService: %s, err: %v", str, err)
	printCmdRes(str, err)
	return err
}

func (this *Service) Status() error {
	log.Println("ServiceManagement.Status not supported")

	log.Printf("config: %s \n", viper.GetString("msg"))
	log.Printf("config: log.logpath%s \n", viper.GetString("log.logpath"))
	log.Printf("config: %v \n", viper.GetStringMap("log")["logpath"])
	log.Printf("config keys: %v \n", viper.AllKeys())
	return nil
}

func (this *Service) StartService() error {
	log.Println("ServiceManager.StartService")

	str, err := this.Start()
	log.Println("StartService: %s, err: %v", str, err)
	printCmdRes(str, err)
	return err
}

func (this *Service) StopService() error {
	log.Println("ServiceManager.StopService")
	str, err := this.Stop()
	log.Println("StopService: %s, err: %v", str, err)
	printCmdRes(str, err)
	return err
}

func (this *Service) PauseService() error {
	log.Println("ServiceManager.PauseServicen not supported ")
	return nil
}

func (this *Service) ContinueService() error {
	log.Println("ServiceManager.ContinueService not supported ")
	return nil
}
