package utils

import (
	"fmt"
	"runtime"
)

func Recover_and_log() {
	if err := recover(); err != nil {

		trace := make([]byte, 1024)
		count := runtime.Stack(trace, true)
		err_msg := fmt.Sprintf("Recover from panic: %s\n", err)
		trace_msg := fmt.Sprintf("Stack of %d bytes: %s\n", count, trace)

		fmt.Println(err_msg)
		fmt.Println(trace_msg)

		// log.Critical(err_msg)
		// log.Critical(trace_msg)
	}
	// log.Flush()
}
