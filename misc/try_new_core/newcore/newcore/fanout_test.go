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

	fset := FanOut(sub, NewStdoutBackend("b1"), NewStdoutBackend("b2"))

	fset_closed_chan := make(chan error)

	time.AfterFunc(time.Second*time.Duration(2), func() {
		// sub will be closed within FanOut
		fset_closed_chan <- fset.Close()
	})

	a := 0
	tick := time.Tick(time.Second * time.Duration(1))
	timeout := time.After(time.Second * time.Duration(6))

main_loop:
	for {
		select {
		case <-fset_closed_chan:
			fmt.Println("TestFanout.fset closed")
		case md, openning := <-sub.Updates():
			if openning == false {
				fmt.Println("TestFanout.sub.Updates() closed")
				break main_loop
			} else {
				fmt.Printf("TestFanout.sub.Updates() still openning: 0x%X\n", &md)
			}
			a += len(*md)
			// t.Log("md: ", md)
		case <-tick:
			a = 0
		case <-timeout:
			t.Error("TestFanout.timed out! something is blocking")
			return
		}
	}
	// t.Error("...")
}

func TestCloseFanoutRepeatly(t *testing.T) {

	// closed_cnt: 999
	// closed_cnt: 1000
	// --- PASS: TestCloseFanoutRepeatly (177.44s)
	// PASS
	// 	ok      github.com/oliveagle/hickwall/misc/try_new_core/newcore/newcore 179.528s
	// tested close 1000 times no bug and with very small amout of memory leaking

	closed_cnt := 0
	expected_closed_cnt := 100

	for cnt := 0; cnt < expected_closed_cnt; cnt++ {
		// a := 0

		fset := FanOut(
			Subscribe(CollectorFactory("c1"), nil),
			NewStdoutBackend("b1"),
			NewStdoutBackend("b2"),
		)

		fset_closed_chan := make(chan error)

		time.AfterFunc(time.Millisecond*time.Duration(10), func() {
			fset_closed_chan <- fset.Close()
		})

		timeout := time.After(time.Second * time.Duration(2))

	close_wait_loop:
		for {
			select {
			case <-fset_closed_chan:
				break close_wait_loop
			case <-timeout:
				t.Error("TestFanout.timed out! something is blocking")
				return
			}
		}
		closed_cnt += 1

		t.Logf("closed_cnt: %d\n", closed_cnt)
	}

	if closed_cnt != expected_closed_cnt {
		t.Error("closed_cnt: %d  !=  expected_closed_cnt: %d", closed_cnt, expected_closed_cnt)
	}

}
