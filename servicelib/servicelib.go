package servicelib

import (
	"fmt"
	"github.com/oliveagle/hickwall/config"
	"github.com/oliveagle/hickwall/daemon"
	"github.com/spf13/viper"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
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

func copyfile(src string, dst string) error {
	//
	in, err := os.Open(src)
	if err != nil {
		log.Printf("Error: cannot open file: %s\n", err)
		os.Exit(2)
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		log.Printf("Error: cannot create file: %s\n", err)
		os.Exit(2)
	}

	_, err = io.Copy(out, in)
	cerr := out.Close()
	if err != nil {
		log.Printf("Error: cannot copy file: %s\n", err)
	}
	return cerr
}

func isUnderRoot(execpath string) bool {
	root := viper.GetString("root")
	return strings.HasPrefix(execpath, root)
}

func getDestPath() string {
	root := viper.GetString("root")
	return root + config.VERSION + config.APP_NAME
}

func copyMeToRoot(dst string) {
	// 如果当前执行文件 不在 root文件夹下, 要先把exec复制到root文件夹下的对应文件
	execPath, err := getExecMyPath()
	if err != nil {
		log.Printf("Error: cannot get executable file path: %v", err)
		os.Exit(2)
	}

	if isUnderRoot(execPath) {
		log.Println("Error: current executable file is already under root: %s", execPath)
		os.Exit(2)
	}

	copyfile(execPath, getDestPath())
}

// type ServiceManager struct{}

func getExecMyPath() (string, error) {
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
