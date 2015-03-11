package main

import (
	"fmt"
	"sync"
	"time"
)

var lock sync.RWMutex
var doing bool

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
