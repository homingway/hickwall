package collectors

import (
	"github.com/oliveagle/hickwall/newcore"
	"testing"
	"time"
)

func TestDummyCollector(t *testing.T) {

	sub := newcore.Subscribe(NewDummyCollector("c1", time.Millisecond*100, 100), nil)

	time.AfterFunc(time.Second*1, func() {
		sub.Close()
	})

	timeout := time.After(time.Second * time.Duration(2))

main_loop:
	for {
		select {
		case md, openning := <-sub.Updates():
			if openning == false {
				break main_loop
			} else {
				t.Log(&md)
			}
		case <-timeout:
			t.Error("timed out! something is blocking")
			break main_loop
		}
	}
}
