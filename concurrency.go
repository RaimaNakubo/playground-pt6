package main

import (
	"fmt"
	"time"
)

func main() {
	go say("world") //calling say() in another goroutine
	say("hello")    //calling say() in current goroutine
}

// function say() prints it's argument five times with 100ms delay
func say(s string) {
	for i := 0; i < 5; i++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Println(s)
	}
}
