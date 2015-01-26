// +build linux darwin

package servicelib

import (
	"github.com/spf13/viper"
	"log"
)

func (this *Service) IsAnInteractiveSession() (bool, error) {
	log.Println("IsAnInteractiveSessioin\r\n")
	// log.Printf("Getegid: %s  \r\n", os.Getegid())
	return false, nil
}

func (this *Service) InstallService() error {
	log.Println("ServiceManager.InstallService\r\n")
	str, err := this.Install()
	log.Println("InstallService: %s, err: %s", str, err)
	return err
}

func (this *Service) RemoveService() error {
	log.Println("ServiceManager.RemoveService\r\n")
	str, err := this.Remove()
	log.Println("RemoveService: %s, err: %s", str, err)
	return err
}

func (this *Service) Status() error {
	log.Println("ServiceManagement.Status not supported\r\n")
	return nil
}

func (this *Service) StartService() error {
	log.Println("ServiceManager.StartService\r\n")

	str, err := this.Start()
	log.Println("StartService: %s, err: %s", str, err)
	return err
}

func (this *Service) StopService() error {
	log.Println("ServiceManager.StopService\r\n")
	str, err := this.Stop()
	log.Println("StopService: %s, err: %s", str, err)
	return err
}

func (this *Service) PauseService() error {
	log.Println("ServiceManager.PauseServicen not supported \r\n")
	return nil
}

func (this *Service) ContinueService() error {
	log.Println("ServiceManager.ContinueService not supported \r\n")
	return nil
}

func (this *Service) Config() error {
	log.Println("Service Config -------")
	// log.Println("config: ", this.config.Name)

	log.Printf("config: %s \n", viper.GetString("msg"))
	log.Printf("config: log.logpath%s \n", viper.GetString("log.logpath"))
	log.Printf("config: %v \n", viper.GetStringMap("log")["logpath"])
	log.Printf("config keys: %v \n", viper.AllKeys())

	return nil
}
