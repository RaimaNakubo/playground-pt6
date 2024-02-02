package main

import (
	"fmt"
	"time"
)

func main() {
	//channels
	s := []int{7, 2, 8, -9, 4, 0} //input data in a slice of int

	c := make(chan int) //creating int channel

	go sum(s[:len(s)/2], c) //summary of the first half of a slice calculates in 2nd goroutine
	go sum(s[len(s)/2:], c) //summary of the second half of a slice calculates in 3rd goroutine

	x, y := <-c, <-c //recieving first and second calculation results from int channel

	//here x recieved result of the second half calculation bc calculation have been finished faster than the first half
	//final calculation happens once both goroutines have completed their's calculations
	fmt.Printf("Summary of: first half - %v, second half - %v, combined - %v\n", x, y, x+y)
	fmt.Println()

	//buffered channels
	ch := make(chan int, 2) //ch is an integer channel with buffer == 2
	ch <- 1
	ch <- 2
	//ch <- 3 fatal error: all goroutines are asleep - deadlock! - goroutine 1 [chan send]
	fmt.Println(<-ch)
	fmt.Println(<-ch)
	//fmt.Println(<-ch) fatal error: all goroutines are asleep - deadlock! - goroutine 1 [chan receive]
	fmt.Println()

	//Range and Close
	ch2 := make(chan int, 10)    //ch2 is an integer channel with buffer == 10
	go fibonacchi(cap(ch2), ch2) //calculating n == capacity of ch2 == buffer numbers of Fibonacchi in a separate goroutine

	for i := range ch2 { //until channel ch2 is closed
		fmt.Println(i) //printing data from channel
	}
	fmt.Println()

	//Select
	ch3 := make(chan int)  //ch3 is an integer channel without buffer; used for transfering calculations
	quit := make(chan int) //quit is an integer channel without buffer; used to stop calculation process

	go func() { //async function call in sub-goroutine
		for i := 0; i < 10; i++ {
			fmt.Println(<-ch3) //which prints 10 values from channel ch3
		}
		quit <- 0 //then sends a signal to stop calculations
	}()

	selectFibonacchi(ch3, quit) //calling calculating function in general goroutine
	fmt.Println()

	//default selection
	tick := time.Tick(100 * time.Millisecond) // tick is a ticking time.Time channel that recieves current time every (100ms) which is
	//(units of time passed in current goroutine)
	boom := time.After(500 * time.Millisecond) // boom is a time.Time channel that recieves current time after (duration) have been elapsed

	for {
		select {
		case <-tick: //if tick happened
			fmt.Println("tick.")
		case <-boom: //if boom happened
			fmt.Println("BOOM!")
			return //loop ends
		default: //if nothing happened
			fmt.Println("nothing..")
			time.Sleep(50 * time.Millisecond) //pause this goroutine for 50ms
		}

	}

}

// function sum summarizes all values inside argument's slice and sends the summary to integer channel
func sum(s []int, c chan int) {
	sum := 0
	for _, v := range s {
		sum += v //calculating summary
	}
	c <- sum //sending summary to the int channel from argument
}

// function fibonacchi() calculates Fibonacchi row for n numbers from argument 1 and send them to int channel from argument 2
func fibonacchi(n int, c chan int) {
	x, y := 0, 1
	for i := 0; i < n; i++ {
		c <- x
		x, y = y, x+y
	}
	close(c) //closing int channel as there would be no more values send
}

// sunction selectFibonacchi() calculates Fibonacchi row for n numbers,
// sending them repetedly to int channel from argument 1,
// until some data appears in int channel from argument 2
func selectFibonacchi(c, quit chan int) {
	x, y := 0, 1
	for {
		select {
		case c <- x:
			x, y = y, x+y
		case <-quit:
			fmt.Println("quit")
			return
		}
	}
}
