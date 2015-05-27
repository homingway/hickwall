package utils

import (
	"strings"
	"testing"
	"time"
)

func TestCommand(t *testing.T) {
	res, err := Command(time.Second, "df")
	t.Log(string(res), err)
	if !strings.Contains(string(res), "Mounted") {
		t.Error("failed")
	}
}
