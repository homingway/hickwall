package collectorlib

import (
	"testing"
	"time"
)

func TestIsDigit(t *testing.T) {
	if IsDigit("1a3") {
		t.Error("1a3: expected false")
	}
	if !IsDigit("029") {
		t.Error("029: expected true")
	}
}

func TestParseInterval_InvalidFormat(t *testing.T) {

	_, err := ParseInterval("1")
	if err == nil {
		t.Error("should raise error")
	}

	_, err = ParseInterval("0m")
	if err == nil {
		t.Error("should raise error 0m")
	}

	_, err = ParseInterval("-1m")
	if err == nil {
		t.Error("should raise error -1m")
	}

	d, err := ParseInterval("1w")
	if err == nil {
		t.Error("should raise error 1w")
	}

	if d != 0 {
		t.Log(d)
		t.Error("-")
	}
}

func TestParseInterval(t *testing.T) {

	d, err := ParseInterval("1s")
	if err != nil {
		t.Error("should not raise error")
	}
	if d == 0 {
		t.Error("d is zero")
	}
	if d != time.Duration(1)*time.Second {
		t.Log(d)
		t.Error("d is not 1s")
	}
}
