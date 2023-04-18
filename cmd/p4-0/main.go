package main

import (
	"fmt"
	"time"
)

func main() {

	// create channel
	ch := make(chan string)

	// function call with goroutine
	go sendData(ch)

	// receive channel data
	fmt.Println("we are here waiting for the value")
	fmt.Println(<-ch)

}

func sendData(ch chan string) {
	time.Sleep(3 * time.Second)
	// data sent to the channel
	//ch <- "Received. Send Operation Successful"
	//fmt.Println("No receiver! Send Operation Blocked")
	close(ch)

}
