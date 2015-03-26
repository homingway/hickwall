package main

// #cgo CPPFLAGS: -I../
// #cgo LDFLAGS: -L../ -lpdh -lstdc++
// #include "wpdh.hpp"
import "C"

import "fmt"

func main() {

	fmt.Println("golang print this line. the rest are from c++\n")

	C.getcpuload()
}
