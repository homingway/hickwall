package main

import (
	"fmt"
	. "github.com/oliveagle/hickwall/misc/tryout_inherit"
	"time"
)

func main() {
	var collectors []Collector

	collectors = append(collectors, &IntervalCollector{Interval: 1, Name: "1"})

	go collectors[0].Run()

	time.Sleep(time.Second * 2)
	collectors[0].SetName("name 1")
	time.Sleep(time.Second * 4)
	// collectors[0].Run()

	collectors = append(collectors, &C_win_pdh{IntervalCollector{Interval: 1, Name: "c_win_pdh"}})
	// collectors = append(collectors, &C_win_pdh{Interval: 1, Name: "c_win_pdh"})
	go collectors[1].Run()

	time.Sleep(time.Second * 2)
	collectors[1].SetName("c_win_pdh name!!")
	time.Sleep(time.Second * 4)
	fmt.Println(collectors[1].GetName())
}
