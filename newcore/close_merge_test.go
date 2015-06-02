package newcore

import (
	"testing"
	"time"
)

func TestCloseMerged(t *testing.T) {
	merged := Merge(
		Subscribe(dummyCollectorFactory("c1"), nil),
		Subscribe(dummyCollectorFactory("c2"), nil),
		Subscribe(dummyCollectorFactory("c3"), nil))

	// Close the subscriptions after some time.
	time.AfterFunc(1*time.Second, func() {
		merged.Close()
	})

	timeout := time.After(time.Duration(2) * time.Second)

	var md MultiDataPoint
	var channel_closed bool

	for {
		select {
		case md, channel_closed = <-merged.Updates():
			if md == nil && channel_closed == false {
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
	// close merged N times
	closed_cnt := 0
	expect_closed_cnt := 100

outer_loop:
	for {

		merged := Merge(
			Subscribe(dummyCollectorFactory("c1"), nil),
			Subscribe(dummyCollectorFactory("c2"), nil),
			Subscribe(dummyCollectorFactory("c3"), nil))

		// Close the subscriptions after some time.
		// time.AfterFunc(1*time.Second, func() {
		time.AfterFunc(10*time.Millisecond, func() {
			merged.Close()
		})

		timeout := time.After(time.Duration(3) * time.Second)

		var md MultiDataPoint
		var channel_closed bool

	inner_loop:
		for {
			if closed_cnt >= expect_closed_cnt {
				break outer_loop
			}
			select {
			case md, channel_closed = <-merged.Updates():
				if md == nil && channel_closed == false {
					t.Log("merged closed")
					closed_cnt += 1
					break inner_loop
				}
			case <-timeout:
				t.Error("merged.Close() timeout")
				break outer_loop
			}
		}
	}
	// t.Error("...")
}
