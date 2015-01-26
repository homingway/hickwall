package servicelib

import (
	"fmt"
	"github.com/takama/daemon"
	"os"
	"path/filepath"
)

type IServiceManager interface {
	InstallService(name string, desc string) error
	RemoveService(name string) error
	StartService(name string) error
	StopService(name string) error
	PauseService(name string) error
	ContinueService(name string) error

	// RunService(name string, isDebug bool)
	Status(name string) string
}

// type ServiceManager struct{}

func exePath() (string, error) {
	prog := os.Args[0]
	p, err := filepath.Abs(prog)
	if err != nil {
		return "", err
	}
	fi, err := os.Stat(p)
	if err == nil {
		if !fi.Mode().IsDir() {
			return p, nil
		}
		err = fmt.Errorf("%s is directory", p)
	}
	if filepath.Ext(p) == "" {
		p += ".exe"
		fi, err := os.Stat(p)
		if err == nil {
			if !fi.Mode().IsDir() {
				return p, nil
			}
			err = fmt.Errorf("%s is directory", p)
		}
	}
	return "", err
}

type Service struct {
	daemon.Daemon
	name string
	desc string
}

func NewService(name, desc string) *Service {
	srv, err := daemon.New(name, desc)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	return &Service{srv, name, desc}
}
