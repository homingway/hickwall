package main

import (
	// "bytes"
	"fmt"
	"time"
)

// func main() {

// 	c := make(chan bool)
// 	// send to a closed channel will panic
// 	close(c)
// 	c <- true
// 	fmt.Println("hahah")

// }
func Once() string {
	time.Sleep(time.Second * time.Duration(30))
	return "from Once"
}

func main() {
	result := make(chan string)

	done := make(chan bool)
	go func() {
		result <- Once()
		done <- true
	}()

	go func() {
		for {
			select {
			case <-done:
				break
			case <-time.After(time.Duration(2) * time.Second):
				result <- "from Timeout"
				break
			}
		}
	}()

	fmt.Println(<-result)
	close(result)

	time.AfterFunc(time.Duration(2), func() {
		fmt.Println("haha")
	})
	panic("haah")
}
