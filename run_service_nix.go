// +build linux darwin

package main

import (
	// "fmt"
	"os"
	"os/signal"
	"syscall"

	log "github.com/oliveagle/hickwall/_third_party/seelog"
	"time"
)

func runService(name string, idDebug bool) (string, error) {
	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	tick := time.Tick(1 * time.Second)

	for {
		select {
		case now := <-tick:
			log.Debug("tick: %v", now)
		case killSignal := <-interrupt:
			if killSignal == os.Interrupt {
				return "Daemon was interruped by system signal", nil
			}
			return "Daemon was killed", nil
		}
	}
}
