package main

import "fmt"

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
