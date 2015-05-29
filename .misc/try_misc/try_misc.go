package main

/*
#include <unistd.h>

void CSleep(unsigned int sec){
    sleep(sec);
}
*/
import "C"

import (
	"fmt"
	// "runtime"
)

func main() {
	fmt.Println("started")
	// runtime.GOMAXPROCS(20)

	for i := 0; i < 5; i++ {
		go func(v int) {
			fmt.Println(v)
		}(i)
	}

	// C.CSleep(5)

	fmt.Println("haahh")
	for {

	}
}
