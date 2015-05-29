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
	sub := Subscribe(dummyCollectorFactory("c1"), nil)

	// fset := FanOut(sub,
	// 	newDummyBackend("b1", time.Second*10),
	// 	newDummyBackend("b2", 0))

	fset := FanOut(sub,
		MustNewDummyBackend("b1", "0", false))

	fset_closed_chan := make(chan error)

	time.AfterFunc(time.Second*time.Duration(1), func() {
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
			break main_loop
		case md, openning := <-sub.Updates():
			if openning == false {
				fmt.Println("TestFanout.sub.Updates() closed")
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
			t.Error("TestFanout.timed out! something is blocking")
			break main_loop
		}
	}
}

func TestFanoutJammingBackend(t *testing.T) {

	// if backend is jamming, we will force close fanout with timeout. and left
	// jamming backend unclosed. so if this process take's too long. there must
	// be some other error happending

	sub := Subscribe(dummyCollectorFactory("c1"), nil)

	fset := FanOut(sub,
		MustNewDummyBackend("b1", "10s", true),
		MustNewDummyBackend("b2", "0", false))

	fset_closed_chan := make(chan error)

	time.AfterFunc(time.Second*time.Duration(1), func() {
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
			break main_loop
		case md, openning := <-sub.Updates():
			if openning == false {
				fmt.Println("TestFanout.sub.Updates() closed")
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
			t.Error("TestFanout.timed out! something is blocking")
			time.Sleep(5) // wait to clean a little bit. test framework will clean and ignore sleep.
			panic("")     // print stack
			break main_loop
		}
	}
}

func TestCloseFanoutRepeatly(t *testing.T) {

	//         fanout_test.go:149: closed_cnt: 997
	//         fanout_test.go:149: closed_cnt: 998
	//         fanout_test.go:149: closed_cnt: 999
	//         fanout_test.go:149: closed_cnt: 1000
	// PASS
	// ok      github.com/oliveagle/hickwall/misc/try_new_core/newcore/newcore 10.204s

	closed_cnt := 0
	expected_closed_cnt := 100

	for cnt := 0; cnt < expected_closed_cnt; cnt++ {
		// a := 0

		fset := FanOut(
			Subscribe(dummyCollectorFactory("c1"), nil),
			MustNewDummyBackend("b1", "0", false),
			MustNewDummyBackend("b2", "0", false),
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
