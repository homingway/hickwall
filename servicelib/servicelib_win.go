// +build windows

package servicelib

import (
	"code.google.com/p/winsvc/eventlog"
	"code.google.com/p/winsvc/mgr"
	"code.google.com/p/winsvc/svc"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"time"
)

func (this *Service) IsAnInteractiveSession() (bool, error) {
	return svc.IsAnInteractiveSession()
}

func (this *Service) StartService() error {
	log.Println("ServiceManager.StartService\r\n")
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	log.Println("Connected mgr\r\n")

	s, err := m.OpenService(this.name)
	if err != nil {
		return fmt.Errorf("could not access service: %v", err)
	}
	defer s.Close()

	log.Println("Opened Service\r\n")

	err = s.Start([]string{"p1", "p2", "p3"})
	if err != nil {
		return fmt.Errorf("could not start service: %v", err)
	}
	log.Println("returned ServiceManager.StartService\r\n")
	return nil
}

func (this *Service) InstallService() error {
	log.Println("ServiceManager.InstallService\r\n")
	exepath, err := exePath()
	if err != nil {
		return err
	}
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(this.name)
	if err == nil {
		s.Close()
		return fmt.Errorf("service %s already exists", this.name)
	}
	s, err = m.CreateService(this.name, exepath, mgr.Config{DisplayName: this.desc})
	if err != nil {
		return err
	}
	defer s.Close()
	err = eventlog.InstallAsEventCreate(this.name, eventlog.Error|eventlog.Warning|eventlog.Info)
	if err != nil {
		s.Delete()
		return fmt.Errorf("SetupEventLogSource() failed: %s", err)
	}
	return nil
}

func (this *Service) RemoveService() error {
	log.Println("ServiceManager.RemoveService\r\n")
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(this.name)
	if err != nil {
		return fmt.Errorf("service %s is not installed", this.name)
	}
	defer s.Close()
	err = s.Delete()
	if err != nil {
		return err
	}
	err = eventlog.Remove(this.name)
	if err != nil {
		return fmt.Errorf("RemoveEventLogSource() failed: %s", err)
	}
	return nil
}

func (this *Service) Status() error {
	log.Println("ServiceManagement.Status --------------------\r\n")

	log.Printf("config: %s \n", viper.GetString("msg"))
	log.Printf("config: log.logpath%s \n", viper.GetString("log.logpath"))
	log.Printf("config: %v \n", viper.GetStringMap("log")["logpath"])
	log.Printf("config keys: %v \n", viper.AllKeys())
	return nil
}

func (this *Service) StopService() error {
	log.Println("ServiceManager.StopService\r\n")
	return controlService(this.name, svc.Stop, svc.Stopped)
}

func (this *Service) PauseService() error {
	log.Println("ServiceManager.PauseService\r\n")
	return controlService(this.name, svc.Pause, svc.Paused)
}

func (this *Service) ContinueService() error {
	log.Println("ServiceManager.ContinueService\r\n")
	return controlService(this.name, svc.Continue, svc.Running)
}

func controlService(name string, c svc.Cmd, to svc.State) error {
	log.Printf("controlService: %s \r\n", name)
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(name)
	if err != nil {
		return fmt.Errorf("could not access service: %v", err)
	}
	defer s.Close()
	status, err := s.Control(c)
	if err != nil {
		return fmt.Errorf("could not send control=%d: %v", c, err)
	}
	timeout := time.Now().Add(10 * time.Second)
	for status.State != to {
		if timeout.Before(time.Now()) {
			return fmt.Errorf("timeout waiting for service to go to state=%d", to)
		}
		time.Sleep(300 * time.Millisecond)
		status, err = s.Query()
		if err != nil {
			return fmt.Errorf("could not retrieve service status: %v", err)
		}
	}
	return nil
}
