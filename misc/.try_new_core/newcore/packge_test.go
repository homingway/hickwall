package newcore

import (
	"fmt"
	"github.com/oliveagle/hickwall/misc/try_new_core/newcore/backends"
	"github.com/oliveagle/hickwall/misc/try_new_core/newcore/collectors"
	"github.com/oliveagle/hickwall/misc/try_new_core/newcore/newcore"
	"testing"
	"time"
)

var (
	_ = fmt.Sprintf("")
	_ = time.Now()
)

func TestMain(t *testing.T) {
	merge := newcore.Merge(
		newcore.Subscribe(collectors.NewDummyCollector("c1", time.Millisecond*100), nil),
		newcore.Subscribe(collectors.NewDummyCollector("c2", time.Millisecond*100), nil),
	)

	// fset := FanOut(merge,
	//  newDummyBackend("b1", time.Second*10),
	//  newDummyBackend("b2", 0))

	fset := newcore.FanOut(merge,
		backends.NewDummyBackend("b1", 0, false),
		backends.NewDummyBackend("b2", 0, false),
	)

	fset_closed_chan := make(chan error)

	time.AfterFunc(time.Second*time.Duration(100), func() {
		// merge will be closed within FanOut
		fset_closed_chan <- fset.Close()
	})

	a := 0
	tick := time.Tick(time.Second * time.Duration(1))
	timeout := time.After(time.Second * time.Duration(106))

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
