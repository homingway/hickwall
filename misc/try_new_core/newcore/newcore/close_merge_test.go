package newcore

import (
	"testing"
	"time"
)

func TestCloseMerged(t *testing.T) {
	merged := Merge(
		Subscribe(CollectorFactory("c1")),
		Subscribe(CollectorFactory("c2")),
		Subscribe(CollectorFactory("c3")))

	// Close the subscriptions after some time.
	time.AfterFunc(1*time.Second, func() {
		merged.Close()
	})

	timeout := time.After(time.Duration(2) * time.Second)

	var dp *DataPoint
	var channel_closed bool

	for {
		select {
		case dp, channel_closed = <-merged.Updates():
			if dp == nil && channel_closed == false {
				t.Log("merged closed")
				return
			}
		case <-timeout:
			t.Error("merged.Close() timeout")
			return
		}
	}
}

func TestCloseMergedInLoop(t *testing.T) {
	done := time.After(time.Duration(20) * time.Second)

outer_loop:
	for {

		merged := Merge(
			Subscribe(CollectorFactory("c1")),
			Subscribe(CollectorFactory("c2")),
			Subscribe(CollectorFactory("c3")))

		// Close the subscriptions after some time.
		time.AfterFunc(1*time.Second, func() {
			merged.Close()
		})

		timeout := time.After(time.Duration(3) * time.Second)

		var dp *DataPoint
		var channel_closed bool

	inner_loop:
		for {
			select {
			case dp, channel_closed = <-merged.Updates():
				if dp == nil && channel_closed == false {
					t.Log("merged closed")
					break inner_loop
				}
			case <-timeout:
				t.Error("merged.Close() timeout")
				break outer_loop
			case <-done:
				break outer_loop
			}
		}
	}
}
