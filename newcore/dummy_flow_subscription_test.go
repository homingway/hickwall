package newcore

import (
	"testing"
	"time"
)

func TestNewDummyFlowSubscription(t *testing.T) {

	sub := NewDummyFlowSubscription("c1", time.Millisecond*100, 1)

	time.AfterFunc(time.Second*1, func() {
		sub.Close()
	})

	timeout := time.After(time.Second * time.Duration(3))

main_loop:
	for {
		select {
		case md, openning := <-sub.Updates():
			if openning == false {
				break main_loop
			} else {
				for _, p := range md {
					t.Log(p)
				}
				// t.Log(&md)
			}
		case <-timeout:
			t.Error("timed out! something is blocking")
			break main_loop
		}
	}
}
