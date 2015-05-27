package utils

import (
	"strings"
	"testing"
	"time"
)

func TestCommand(t *testing.T) {
	res, err := Command(time.Second, "cmd.exe", "/k", "dir", "c:\\")
	t.Log(string(res), err)
	if !strings.Contains(string(res), "<DIR>") {
		t.Error("failed")
	}
}
