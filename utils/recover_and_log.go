package utils

import (
	"fmt"
	"runtime"
)

var trace = make([]byte, 1024, 1024)

func Recover_and_log() {
	if err := recover(); err != nil {
		count := runtime.Stack(trace, true)
		err_msg := fmt.Sprintf("Recover from panic: %s\n", err)
		trace_msg := fmt.Sprintf("Stack of %d bytes: %s\n", count, trace)
		fmt.Println(err_msg)
		fmt.Println(trace_msg)
	}
}
