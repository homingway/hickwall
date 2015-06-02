package utils

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"os/exec"
	"time"

	"log"
)

var (
	// ErrPath is returned by Command if the program is not in the PATH.
	ErrPath = errors.New("program not in PATH")
	// ErrTimeout is returned by Command if the program timed out.
	ErrTimeout = errors.New("program killed after timeout")
)

// Command executes the named program with the given arguments. If it does not
// exit within timeout, it is sent SIGINT (if supported by Go). After
// another timeout, it is killed.
func Command(timeout time.Duration, name string, arg ...string) ([]byte, error) {
	// avoid leaking param
	var command_name = name[:]
	var args []string
	for _, a := range arg {
		args = append(args, a)
	}

	if _, err := exec.LookPath(command_name); err != nil {
		return nil, ErrPath
	}
	log.Printf("executing command: %v %v\n", command_name, args)
	c := exec.Command(command_name, args...)
	var b bytes.Buffer
	c.Stdout = &b
	done := make(chan error, 1)

	// this cmd escapes to heap. but better than leaking.
	go func(cmd exec.Cmd) {
		done <- cmd.Run()
	}(*c)

	interrupt := time.After(timeout)
	kill := time.After(timeout * 2)
	for {
		select {
		case err := <-done:
			return b.Bytes(), err
		case <-interrupt:
			c.Process.Signal(os.Interrupt)
		case <-kill:
			// todo: figure out if this can leave the done chan hanging open
			c.Process.Kill()
			return nil, ErrTimeout
		}
	}
}

// ReadCommand runs command name with args and calls line for each line from its
// stdout. Command is interrupted (if supported by Go) after 10 seconds and
// killed after 20 seconds.
func ReadCommand(line func(string) error, name string, arg ...string) error {
	return ReadCommandTimeout(time.Second*10, line, name, arg...)
}

// ReadCommandTimeout is the same as ReadCommand with a specifiable timeout.
func ReadCommandTimeout(timeout time.Duration, line func(string) error, name string, arg ...string) error {
	b, err := Command(timeout, name[:], arg...)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(bytes.NewBuffer(b))
	for scanner.Scan() {
		if err := line(scanner.Text()); err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		log.Printf("ERROR: %v: %v\n", name[:], err)
	}
	return nil
}
