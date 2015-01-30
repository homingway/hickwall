// +build windows

package daemon

import (
	"os"
	"os/exec"
	"os/user"
	"strings"
)

// Service constants
const (
	rootPrivileges = "You must have root user privileges. Possibly using 'sudo' command should help"
	success        = "\t\t\t[  OK  ]" // Show colored "OK"
	failed         = "\t\t\t[FAILED]" // Show colored "FAILED"
)

func IsExecutable(path string) (bool, error) {
	in, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer in.Close()

	stat, err := in.Stat()
	if err != nil {
		return false, err
	}

	if stat.Mode().IsRegular() && strings.HasSuffix(path, ".exe") {
		return true, nil
	} else {
		return false, nil
	}
}

// Lookup path for executable file
func executablePath(name string) (string, error) {
	if path, err := exec.LookPath(name); err == nil {
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			return execPath()
		}
		return path, nil
	}
	return execPath()
}

// Check root rights to use system service
func checkPrivileges() bool {

	if user, err := user.Current(); err == nil && user.Gid == "0" {
		return true
	}
	return false
}
