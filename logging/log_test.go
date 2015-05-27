package logging

import (
	"testing"
)

func TestLogPerf(t *testing.T) {
	log := MustGetLogger()

	for i := 0; i < 100; i++ {
		log.Info("this is msg %s ok", "haha")
		log.Error("haha")
		log.Critical("hahah")
	}

	log.Info(GetModule())
}

func TestGetModule(t *testing.T) {
	mod := GetModule()
	if mod != "github.com/oliveagle/hickwall/logging" {
		t.Error("GetModule")
	}
}
