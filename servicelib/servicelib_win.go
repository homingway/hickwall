// +build windows

package servicelib

import (
	"code.google.com/p/winsvc/eventlog"
	"code.google.com/p/winsvc/mgr"
	"code.google.com/p/winsvc/svc"
	"fmt"
	// log "github.com/cihub/seelog"
	// "github.com/spf13/viper"
	"time"
)

func IsAnInteractiveSession() (bool, error) {
	return svc.IsAnInteractiveSession()
}

func (this *Service) StartService(args ...string) error {
	// log.Debug("ServiceManager.StartService\r\n")
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	// log.Debug("Connected mgr\r\n")

	s, err := m.OpenService(this.name)
	if err != nil {
		return fmt.Errorf("could not access service: %v", err)
	}
	defer s.Close()

	// log.Debug("Opened Service\r\n")

	// err = s.Start([]string{"p1", "p2", "p3"})
	err = s.Start(args)
	if err != nil {
		return fmt.Errorf("could not start service: %v", err)
	}
	return nil
}

func (this *Service) InstallService() error {
	// log.Debug("ServiceManager.InstallService")

	exepath, err := exePath()
	if err != nil {
		// log.Error("exePath error", err)
		return err
	}
	m, err := mgr.Connect()
	if err != nil {
		// log.Error("cannot connect to windows service manager.", err)
		return err
	}
	defer m.Disconnect()

	// log.Info("this.name: ", this.name)
	s, err := m.OpenService(this.name)
	if err == nil {
		s.Close()
		// log.Errorf("service %s already exists", this.name)
		return fmt.Errorf("service %s already exists", this.name)
	}

	s, err = m.CreateService(this.name, exepath, mgr.Config{DisplayName: this.desc})
	if err != nil {
		// log.Error("Failed to create service: ", this.name, exepath, this.desc, err)
		return err
	}
	defer s.Close()

	err = eventlog.InstallAsEventCreate(this.name, eventlog.Error|eventlog.Warning|eventlog.Info)
	if err != nil {
		s.Delete()
		// log.Errorf("SetupEventLogSource() failed: %s", err)
		return fmt.Errorf("SetupEventLogSource() failed: %s", err)
	}
	return nil
}

func (this *Service) RemoveService() error {
	// log.Debug("ServiceManager.RemoveService\r\n")

	m, err := mgr.Connect()
	if err != nil {
		// log.Error("cannot connect to windows service manager.", err)
		return err
	}
	defer m.Disconnect()

	s, err := m.OpenService(this.name)
	if err != nil {
		// log.Errorf("service %s is not installed", this.name)
		return fmt.Errorf("service %s is not installed", this.name)
	}
	defer s.Close()

	err = s.Delete()
	if err != nil {
		return err
	}
	err = eventlog.Remove(this.name)
	if err != nil {
		// log.Errorf("RemoveEventLogSource() failed: %s", err)
		return fmt.Errorf("RemoveEventLogSource() failed: %s", err)
	}
	return nil
}

func (this *Service) Status() (State, error) {
	// log.Debug("ServiceManager.Status\r\n")

	m, err := mgr.Connect()
	if err != nil {
		// log.Error("cannot connect to windows service manager.", err)
		return State(Unknown), err
	}
	defer m.Disconnect()

	s, err := m.OpenService(this.name)
	if err != nil {
		// log.Errorf("service %s is not installed", this.name)
		return State(Unknown), fmt.Errorf("service %s is not installed", this.name)
	}
	defer s.Close()

	status, err := s.Query()
	if err != nil {
		// log.Error("cannot query service status: ", err)
		return State(Unknown), fmt.Errorf("cannot query service status", err)
	}

	switch status.State {
	case svc.Stopped:
		return State(Stopped), nil
	case svc.StartPending:
		return State(StartPending), nil
	case svc.StopPending:
		return State(StopPending), nil
	case svc.Running:
		return State(Running), nil
	case svc.ContinuePending:
		return State(ContinuePending), nil
	case svc.PausePending:
		return State(PausePending), nil
	case svc.Paused:
		return State(Paused), nil
	default:
		// log.Errorf("unknown service state: svc: %s, state: %v", this.name, status.State)
		return State(Unknown), fmt.Errorf("unknown state")
	}
}

func (this *Service) StopService() error {
	// log.Debug("ServiceManager.StopService\r\n")
	return controlService(this.name, svc.Stop, svc.Stopped)
}

func (this *Service) PauseService() error {
	// log.Debug("ServiceManager.PauseService\r\n")
	return controlService(this.name, svc.Pause, svc.Paused)
}

func (this *Service) ContinueService() error {
	// log.Debug("ServiceManager.ContinueService\r\n")
	return controlService(this.name, svc.Continue, svc.Running)
}

func controlService(name string, c svc.Cmd, to svc.State) error {
	// log.Debug("controlService: %s \r\n", name)

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
