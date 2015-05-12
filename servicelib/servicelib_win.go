// +build windows

package servicelib

import (
	"code.google.com/p/winsvc/eventlog"
	"code.google.com/p/winsvc/mgr"
	"code.google.com/p/winsvc/svc"
	"fmt"
	"time"
)

func IsAnInteractiveSession() (bool, error) {
	return svc.IsAnInteractiveSession()
}

func (this *Service) InstallService() error {
	// open logger
	elog, err := eventlog.Open(this.name)
	if err != nil {
		// install logger
		err = eventlog.InstallAsEventCreate(this.name, eventlog.Error|eventlog.Warning|eventlog.Info)
		if err != nil {
			return fmt.Errorf("SetupEventLogSource() failed: %s", err)
		}

		// open
		elog, err = eventlog.Open(this.name)
		if err != nil {
			return fmt.Errorf("Cannot open eventlog")
		}

	}
	defer elog.Close()

	m, err := mgr.Connect()
	if err != nil {
		elog.Error(3, fmt.Sprintf("install service: cannot connect to windows service manager.", err))
		return err
	}
	defer m.Disconnect()

	s, err := m.OpenService(this.name)
	if err == nil {
		s.Close()
		elog.Error(3, fmt.Sprintf("install service: service %s already exists", this.name))
		return fmt.Errorf("service %s already exists", this.name)
	}

	s, err = m.CreateService(this.name, this.path, mgr.Config{DisplayName: this.desc})
	if err != nil {
		elog.Error(3, fmt.Sprintf("install service: Failed to create service: %v", err))
		return err
	}
	defer s.Close()

	err = this.ConfigServiceAutoStart()
	if err != nil {
		elog.Error(3, fmt.Sprintf("install service: cannot set service to auto start: %v", err))
		return err
	}
	return nil
}

func (this *Service) RemoveService() error {
	l, err := eventlog.Open(this.name)
	if err != nil {
		return fmt.Errorf("Cannot open eventlog")
	}

	m, err := mgr.Connect()
	if err != nil {
		l.Error(3, fmt.Sprintf("remove service: cannot connect to windows service manager.", err))
		return err
	}
	defer m.Disconnect()

	s, err := m.OpenService(this.name)
	if err != nil {
		l.Error(3, fmt.Sprintf("remove service: service %s doesn't exists", this.name))
		return fmt.Errorf("service %s doesn't exists", this.name)
	}

	err = s.Delete()
	if err != nil {
		return err
	}

	l.Info(1, "Service Removed. removing event log")
	l.Close()

	eventlog.Remove(this.name)
	return nil
}

func (this *Service) StartService(args ...string) error {
	_, err := queryService(this.name, func(s *mgr.Service, l *eventlog.Log) (State, error) {
		l.Info(1, "Starting Service")
		// tested and checked svc source code, args doesn't work. don't know why.
		err := s.Start(args)
		if err != nil {
			l.Error(3, fmt.Sprintf("could not start service: %v", err))
			return State(Unknown), fmt.Errorf("could not start service: %v", err)
		}

		l.Info(1, "Service Started")
		return State(Unknown), nil
	})
	return err
}

func (this *Service) Status() (State, error) {
	return queryService(this.name, func(s *mgr.Service, l *eventlog.Log) (State, error) {

		status, err := s.Query()
		if err != nil {
			return State(Unknown), fmt.Errorf("cannot queryService service status", err)
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
	})
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

func (this *Service) ConfigServiceAutoStart() error {
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

	c, err := s.Config()
	if err != nil {
		return fmt.Errorf("could not get service config: %v", err)
	}

	// make service automatically start
	c.StartType = mgr.StartAutomatic
	err = s.UpdateConfig(c)
	if err != nil {
		return fmt.Errorf("could not update service config: %v", err)
	}
	return nil
}

func queryService(name string, fn func(s *mgr.Service, l *eventlog.Log) (State, error)) (State, error) {
	l, err := eventlog.Open(name)
	if err != nil {
		return State(Unknown), fmt.Errorf("Cannot open eventlog")
	}
	defer l.Close()

	m, err := mgr.Connect()
	if err != nil {
		l.Error(3, fmt.Sprintf("cannot connect to mgr: %v", err))
		return State(Unknown), err
	}
	defer m.Disconnect()

	s, err := m.OpenService(name)
	if err != nil {
		l.Error(3, fmt.Sprintf("cannot open service: %v", err))
		return State(Unknown), fmt.Errorf("service %s is not installed", name)
	}
	defer s.Close()

	return fn(s, l)
}

func controlService(name string, c svc.Cmd, to svc.State) error {
	// log.Debug("controlService: %s \r\n", name)

	_, err := queryService(name, func(s *mgr.Service, l *eventlog.Log) (State, error) {

		status, err := s.Control(c)
		if err != nil {
			return State(Unknown), fmt.Errorf("could not send control=%d: %v", c, err)
		}

		timeout := time.Now().Add(10 * time.Second)

		for status.State != to {
			if timeout.Before(time.Now()) {
				return State(Unknown), fmt.Errorf("timeout waiting for service to go to state=%d", to)
			}
			time.Sleep(300 * time.Millisecond)
			status, err = s.Query()
			if err != nil {
				return State(Unknown), fmt.Errorf("could not retrieve service status: %v", err)
			}
		}
		return State(Unknown), nil
	})
	return err

}
