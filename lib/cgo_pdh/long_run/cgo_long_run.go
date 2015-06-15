package main

import (
	"fmt"
	"github.com/oliveagle/hickwall/lib/cgo_pdh"
	"io"
	"log"
	"os"
	"time"
)

func main() {
	file, err := os.OpenFile("d:\\cgo_long_run.log", os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("cannot open d:\\cgo_long_run.log for write", err)
		return
	}
	log.SetOutput(io.MultiWriter(os.Stdout, file))
	log.SetFlags(log.Ldate | log.Ltime)
	log.Println("Started - run 100 hours")

	pc := cgo_pdh.NewPdhCollector()
	defer pc.Close()
	pc.AddCounter("\\Process(long_run)\\Working Set - Private")
	data, _ := pc.CollectAllDouble()

	tick := time.Tick(time.Second * 1)
	tickClose := time.Tick(time.Second * 2)
	done := time.After(time.Hour * 100)

	first_value := 0
	last_value := 0
	delta := 0

	for {
		select {
		case <-tickClose:
			pc.Close()
			pc = cgo_pdh.NewPdhCollector()
			pc.AddCounter("\\Process(long_run)\\Working Set - Private")
			log.Println("close and recreate pdh collector")
		case <-tick:
			data, err = pc.CollectAllDouble()
			for _, d := range data {
				if first_value == 0 {
					first_value = int(d) / 1024
				}
				last_value = int(d) / 1024
				delta = last_value - first_value
				log.Printf("first: %d, last: %d, delta: %d\n", first_value, last_value, delta)
			}
		case <-done:
			return
		}
	}
	log.Println("Finished")
}
