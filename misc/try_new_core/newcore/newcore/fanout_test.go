package newcore

import (
	"fmt"
	"testing"
	"time"
)

var (
	_ = fmt.Sprintf("")
	_ = time.Now()
)

func TestFanout(t *testing.T) {
	sub := Subscribe(CollectorFactory("c1"), nil)
	// sub

	// FanOut(sub, ...)

	// fset := FanOut(sub, NewStdoutBackend("b1"))
	fset := FanOut(sub, NewStdoutBackend("b1"))
	time.AfterFunc(time.Second*time.Duration(3), func() {
		// panic("")
		// t.Error("...")
		// return
		fset.Close()
		sub.Close()

	})

	go fset.Run()

	a := 0
	tick := time.Tick(time.Second * time.Duration(1))
	done := time.After(time.Second * time.Duration(3))

main_loop:
	for {
		select {
		case md, closed := <-sub.Updates():
			if closed == false {
				break main_loop
			}
			a += len(*md)
			// t.Log("md: ", md)
		case <-tick:
			a = 0
		case <-done:
			t.Error("...")
			return
		}
	}
	t.Error("...")
}
