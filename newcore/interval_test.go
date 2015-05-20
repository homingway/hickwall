package newcore

import (
	"testing"
	"time"
)

func TestNewInterval(t *testing.T) {
	intv := NewInterval("1s")
	d := intv.MustDuration(time.Second * 10)
	if d != time.Second {
		t.Error("...")
	}
}

func TestIntervalMustDuration(t *testing.T) {
	var i1 Interval
	i1 = "1s"
	t.Log(i1.MustDuration(time.Second * 10))

	if i1.MustDuration(time.Second*10) != time.Second {
		t.Error("failed")
	}

	i1 = "0"
	t.Log(i1.MustDuration(time.Second * 10))

	if i1.MustDuration(time.Second*10) != time.Duration(0) {
		t.Error("failed")
	}

	i1 = "-1s"
	t.Log(i1.MustDuration(time.Second * 10))

	if i1.MustDuration(time.Second*10) != time.Second*10 {
		t.Error("failed")
	}

	i1 = ""
	t.Log(i1.MustDuration(time.Second * 10))

	if i1.MustDuration(time.Second*10) != time.Second*10 {
		t.Error("failed")
	}
}
