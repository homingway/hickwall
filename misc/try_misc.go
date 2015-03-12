package main

import (
	"fmt"
	"sync"
	"time"

	"regexp"
	"strings"
)

var lock sync.RWMutex
var doing bool

var pat_influxdb_version = regexp.MustCompile(`[v]?\d+\.\d+\.\d+[\S]*`)

func do(n int) {
	fmt.Println(" ! do ", n)
	if doing == true {
		return
	}

	lock.Lock()
	defer lock.Unlock()

	doing = true
	fmt.Println("  --->> doing ", n)
	time.Sleep(time.Millisecond * 2000)
	doing = false
}

func main() {

	// a := " InfluxDB v0.8.8 (git: afde71e) (leveldb: 1.15)"
	// a := "ver: 0.9.0-rc7 asdfasdfadsv 1234sdf"
	a := "0.9.0-rc7"
	ss := pat_influxdb_version.FindAllString(a, -1)
	version := ""
	if len(ss) > 0 {
		version = ss[0]
		if strings.HasPrefix(version, "v") {
			version = strings.TrimLeft(version, "v")
			// fmt.Printf("*%s*\n", version)
		}
	}
	fmt.Printf("*%s*\n", version)

	return

	tick := time.Tick(time.Millisecond * 1000)
	done := time.Tick(time.Second * 5)

	cnt := 0
loop:
	for {
		select {
		case <-tick:
			cnt += 1
			fmt.Println(" <- tick ----------------------", cnt, len(tick))
			go do(cnt)
		case <-done:
			fmt.Println(" <- done --------------------- done -")
			break loop
		}
	}
}
