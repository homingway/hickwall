package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	tick := time.Tick(time.Second * time.Duration(1))
	done := time.Tick(time.Second * time.Duration(10))

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGKILL, syscall.SIGTERM)

loop:
	for {
		select {
		case <-tick:
			fmt.Println(".")
		case s := <-c:
			fmt.Println("Got Signal: ", s)
		case <-done:
			break loop
		}
	}
}
