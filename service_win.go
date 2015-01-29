// +build windows

package main

import (
	"code.google.com/p/winsvc/svc"
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/spf13/viper"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// BUG(brainman): MessageBeep Windows api is broken on Windows 7,
// so this example does not beep when runs as service on Windows 7.

var (
	beepFunc = syscall.MustLoadDLL("user32.dll").MustFindProc("MessageBeep")
)

func beep() {
	log.Info("beep")
	beepFunc.Call(0xffffffff)
}

func acceptConnection(listener net.Listener, listen chan<- net.Conn) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		listen <- conn
	}
}

func handleClient(client net.Conn) {
	for {
		buf := make([]byte, 4096)
		numbytes, err := client.Read(buf)
		fmt.Printf("numbytes: %d, err: %s, buf: %v \r\n", numbytes, err, buf[:numbytes])
		if numbytes == 0 || err != nil {
			// EOF, close connection
			return
		}
		if numbytes == 2 && buf[0] == 13 && buf[1] == 10 {
			// [13 10]  "\r\n"
		} else {
			now := time.Now()
			str := fmt.Sprintf("%s: %s\r\n", now.Local().Format("15:04:05.999999999"), buf)
			client.Write([]byte(str))
		}
	}
}

type myservice struct{}

func (this *myservice) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	log.Info("myservice.Execute\r\n")
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown | svc.AcceptPauseAndContinue
	changes <- svc.Status{State: svc.StartPending}

	fasttick := time.Tick(500 * time.Millisecond)
	slowtick := time.Tick(2 * time.Second)
	tick := fasttick

	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

	// execute Run
	go serveConn(args, r, changes)

	// major loop for signal processing.
loop:
	for {
		select {
		case <-tick:
			beep()
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
				// testing deadlock from https://code.google.com/p/winsvc/issues/detail?id=4
				time.Sleep(100 * time.Millisecond)
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				break loop
			case svc.Pause:
				changes <- svc.Status{State: svc.Paused, Accepts: cmdsAccepted}
				tick = slowtick
			case svc.Continue:
				changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
				tick = fasttick
			default:
				log.Error("unexpected control request #%d", c)
			}
		}
	}
	changes <- svc.Status{State: svc.StopPending}
	return
}

func runService(name string, isDebug bool) {
	log.Debug("runService: starting %s service \r\n", name)
	err := svc.Run(name, &myservice{})
	if err != nil {
		log.Debug("runService: Error: %s service failed: %v\r\n", name, err)
		return
	}
	log.Debug("runService: %s service stopped\r\n", name)
}

func serveConn(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (bool, error) {
	log.Println("serveConn\r\n")

	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Set up listener for defined host and port
	port := viper.GetString("port")
	listener, err := net.Listen("tcp", port)
	if err != nil {
		return false, err
	}

	// set up channel on which to send accepted connections
	listen := make(chan net.Conn, 100)
	go acceptConnection(listener, listen)

	// loop work cycle with accept connections or interrupt
	// by system signal
	log.Println("Manage() loop\r\n")
	for {
		select {
		case conn := <-listen:
			go handleClient(conn)
		case killSignal := <-interrupt:
			log.Println("Got signal:", killSignal, "\r\n")
			log.Println("Stoping listening on ", listener.Addr(), "\r\n")
			listener.Close()
			if killSignal == os.Interrupt {
				return false, fmt.Errorf("Daemon was interruped by system signal")
			}
			return false, fmt.Errorf("Daemon was killed")
		}
	}
	return true, nil
}
