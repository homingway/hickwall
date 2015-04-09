package main

import (
	"fmt"
	"github.com/oliveagle/hickwall/backends"
	"github.com/oliveagle/hickwall/collectorlib"
	. "github.com/oliveagle/hickwall/collectors"
	"github.com/oliveagle/hickwall/config"
	"time"
)

func main() {

	config.Conf.Transport_stdout.Enabled = true
	// config.Conf.Transport_stdout.Enabled = false

	stdout, _ := backends.GetBackendByName("stdout")
	go stdout.Run()

	cs := GetBuiltinCollectorByName("builtin_hickwall_client")
	if cs != nil {

		fmt.Println(" ++ builtin_collectors: ", cs)

		ch := make(chan collectorlib.MultiDataPoint)

		fmt.Println("Enabled: ", cs.Enabled())

		go cs.Run(ch)

		done := time.After(time.Second * 3)
	loop:
		for {
			select {
			case md, _ := <-ch:
				stdout.Write(md)
			case <-done:
				break loop
			}
		}
	}
}
