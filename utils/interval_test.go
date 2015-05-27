package utils

import (
	"testing"
	"time"
)

func TestMustPositiveIntervalEmptyDefault(t *testing.T) {
	d := MustPositiveInterval("", time.Second)
	t.Log(d)
	if d != time.Second {
		t.Error("...")
	}

	d = MustPositiveInterval("100", time.Second)
	t.Log(d)
	if d != time.Second {
		t.Error("...")
	}
}

func TestMustPositiveIntervalZero(t *testing.T) {
	d := MustPositiveInterval("0", time.Second)
	t.Log(d)
	if d != time.Duration(0) {
		t.Error("...")
	}

	d = MustPositiveInterval("0ms", time.Second)
	t.Log(d)
	if d != time.Duration(0) {
		t.Error("...")
	}
}

func TestMustPositiveIntervalSuccess(t *testing.T) {
	d := MustPositiveInterval("100ms", time.Second)
	t.Log(d)
	if d != time.Millisecond*100 {
		t.Error("...")
	}
}

func TestMustPositiveIntervalNagRetDefault(t *testing.T) {
	d := MustPositiveInterval("-100ms", time.Second)
	t.Log(d)
	if d != time.Second {
		t.Error("...")
	}
}
