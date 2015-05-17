package newcore

import (
	"testing"
)

func TestFanout(t *testing.T) {
	sub := Subscribe(CollectorFactory("c1"), nil)
	// sub

	// FanOut(sub, ...)
}
