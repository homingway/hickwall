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

func TestHookBackend(t *testing.T) {
	bk := NewHookBackend()
	fset := FanOut(Subscribe(dummyCollectorFactory("c1"), nil), bk)

	fset_closed_chan := make(chan error)

	time.AfterFunc(time.Second*time.Duration(1), func() {
		// merge will be closed within FanOut
		fset_closed_chan <- fset.Close()
	})

	a := 0
	tick := time.Tick(time.Second * time.Duration(1))
	timeout := time.After(time.Second * time.Duration(3))

main_loop:
	for {
		select {
		case <-fset_closed_chan:
			fmt.Println("fset closed")
			break main_loop
		case md, openning := <-bk.Hook():
			if openning == false {
				fmt.Println("HookBackend closed")
				break main_loop
			} else {
				fmt.Printf(".")
			}
			a += len(md)
			//			t.Log("md: ", md)
		case <-tick:
			a = 0
		case <-timeout:
			t.Error("timed out! something is blocking")
			break main_loop
		}
	}
}
