package utils

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
)

func HttpPprofServe(port int) {
	go func() {
		log.Println(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), nil))
	}()
}
