package newcore

import (
	"math"
	"testing"
	"time"
)

func TestGetHostname(t *testing.T) {
	h := GetHostname()

	t.Log(h)

	if h == "" {
		t.Error("...")
	}
}

func equal_float64(a, b float64, percision float64) bool {
	diff := math.Abs(a - b)
	if a == b {
		return true
	} else {
		if diff <= math.Abs(percision) {
			return true
		} else {
			return false
		}
	}
}

func TestNow(t *testing.T) {
	var (
		now1, now2, now3 time.Time
	)

	now1 = Now()
	now2 = Now()
	t.Log(now1)
	t.Log(now2)
	if now1 != now2 {
		t.Error("now1 and now2 are differnt: ", now1, now2)
	}

	// first_tick := time.After(time.Millisecond * 1000)
	second_tick := time.After(time.Millisecond * 1000)

	for {
		select {
		// case <-first_tick:

		case <-second_tick:
			now3 = Now()
			t.Log(now3)
			s := now3.Sub(now2)
			t.Log(s, s.Seconds())

			if !equal_float64(s.Seconds(), 1, 0.01) {
				t.Error("tick is not 1 second", now3.Sub(now2))
			}

			return
		}
	}

}
