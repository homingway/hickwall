// +build linux darwin

package main

import (
	"fmt"
	// "log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/viper"
	"time"
)

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

func runService(name string, idDebug bool) (string, error) {
	// log.Println("runService()\r\n")

	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Set up listener for defined host and port
	port := viper.GetString("port")
	// log.Printf("port: %v\n", port)

	listener, err := net.Listen("tcp", port)
	if err != nil {
		return "Possibly was a problem with the port binding", err
	}

	// set up channel on which to send accepted connections
	listen := make(chan net.Conn, 100)
	go acceptConnection(listener, listen)

	// loop work cycle with accept connections or interrupt
	// by system signal
	// log.Println("Manage() loop\r\n")
	for {
		select {
		case conn := <-listen:
			go handleClient(conn)
		case killSignal := <-interrupt:
			// log.Println("Got signal:", killSignal, "\r\n")
			// log.Println("Stoping listening on ", listener.Addr(), "\r\n")
			listener.Close()
			if killSignal == os.Interrupt {
				return "Daemon was interruped by system signal", nil
			}
			return "Daemon was killed", nil
		}
	}
}
