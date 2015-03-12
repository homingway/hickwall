package utils

import (
	"testing"
)

func Test_SMA_1(t *testing.T) {
	m := SMA{N: 3}
	for i := 0.0; i < 5.0; i += 0.3 {
		r, e := m.Calc(float64(i))
		t.Logf("%0.5f, %v", r, e)
	}
	// t.Error("--")
}

func Test_SMA_2(t *testing.T) {
	m := SMA{N: 3}
	r, e := m.Calc(0.3)
	if e == nil {
		t.Error("Raise Error")
	}
	r, e = m.Calc(0.3)
	if e == nil {
		t.Error("Raise Error")
	}
	r, e = m.Calc(0.3)
	if e != nil {
		t.Error("Should not Raise Error", e)
	}
	if r != 0.3 {
		t.Error("Wrong: should be 0.3", r)
	}
}
