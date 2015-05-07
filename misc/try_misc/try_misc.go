package main

import (
	"fmt"
	"os"
)

func main() {
	h, _ := os.Hostname()
	fmt.Println(h)
}
