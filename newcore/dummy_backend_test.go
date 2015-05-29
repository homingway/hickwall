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

func TestDummyBackend(t *testing.T) {
	merge := Merge(
		Subscribe(NewDummyCollector("c1", time.Millisecond*100, 100), nil),
		Subscribe(NewDummyCollector("c2", time.Millisecond*100, 100), nil),
	)

	fset := FanOut(merge,
		MustNewDummyBackend("b1", "0", false),
		MustNewDummyBackend("b2", "0", false),
	)

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
		case md, openning := <-merge.Updates():
			if openning == false {
				fmt.Println("merge.Updates() closed")
				break main_loop
			} else {
				fmt.Printf(".")
				// fmt.Printf("TestFanout.sub.Updates() still openning: 0x%X\n", &md)
			}
			a += len(*md)
			// t.Log("md: ", md)
		case <-tick:
			a = 0
		case <-timeout:
			t.Error("timed out! something is blocking")
			break main_loop
		}
	}
}
