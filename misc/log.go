package main

import (
	"github.com/oliveagle/hickwall/logging"
)

func main() {
	log := logging.MustGetLogger("name")

	for i := 0; i < 100; i++ {
		log.Info("this is msg %s ok", "haha")
		log.Error("haha")
		log.Critical("hahah")
	}

	log.Info(logging.GetCurrentDir())
}
