// +build linux darwin

package main

import (
	// "fmt"
	"os"
	"os/signal"
	"syscall"

	log "github.com/cihub/seelog"

	"github.com/oliveagle/go-metrics"
	"net"
	"time"
)

func runService(idDebug bool) (string, error) {
	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	tick := time.Tick(1 * time.Second)

	// go service_process()
	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:2003")
	go metrics.Graphite(metrics.DefaultRegistry, 1*time.Second, "metrics", addr)

	c := metrics.NewCounter()
	metrics.Register("foo", c)

	for {
		select {
		case now := <-tick:
			c.Inc(47)
			log.Trace("tick: %v", now)

		case killSignal := <-interrupt:
			if killSignal == os.Interrupt {
				return "Daemon was interruped by system signal", nil
			}
			return "Daemon was killed", nil
		}
	}
}
