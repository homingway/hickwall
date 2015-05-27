package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	files, err := ioutil.ReadDir("./a")
	if err != nil {
		fmt.Println("err: ", err)
	}
	for _, f := range files {
		fmt.Println(f.Name())
	}
}
