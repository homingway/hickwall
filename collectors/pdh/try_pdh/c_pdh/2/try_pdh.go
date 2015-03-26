package main

/*
#cgo CFLAGS: -I../..
#cgo LDFLAGS: -L../.. -lpdh

#include "try_pdh.h"
*/
import "C"

import "fmt"

func main() {

	fmt.Println("this is built with `go build` ")
	C.getcpuload()
}
