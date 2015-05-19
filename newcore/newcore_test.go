package newcore

import (
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

func TestNow(t *testing.T) {
	now1 := Now()
	now2 := Now()
	now3 := Now()
	time.AfterFunc(time.Second*1, func() {
		now3 = Now()
	})

	time.Sleep(time.Second * 1)
	t.Log(now1)
	t.Log(now2)
	t.Log(now3)
	if now1 != now2 {
		t.Error("now1 and now2 are differnt: ", now1, now2)
	}

	s := now3.Sub(now2)
	t.Log(s, s.Seconds())

	if s.Seconds() != 1 {
		t.Error("tick is not 1 second", now3.Sub(now2))
	}
}
