package main

import (
	"bytes"
	"fmt"
)

func main() {
	buf := bytes.NewBuffer(make([]byte, 0, 10))

	for i := 0; i < 100; i++ {
		fmt.Fprintf(buf, "hahah %s", "hello")
	}

	fmt.Println(buf.String())
	buf.Reset()

}
