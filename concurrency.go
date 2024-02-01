package main

import "fmt"

func main() {
	s := []int{7, 2, 8, -9, 4, 0} //input data in a slice of int

	c := make(chan int) //creating int channel

	go sum(s[:len(s)/2], c) //summary of the first half of a slice calculates in 2nd goroutine
	go sum(s[len(s)/2:], c) //summary of the second half of a slice calculates in 3rd goroutine

	x, y := <-c, <-c //recieving first and second calculation results from int channel

	//here x recieved result of the second half calculation bc calculation have been finished faster than the first half
	//final calculation happens once both goroutines have completed their's calculations
	fmt.Printf("Summary of: first half - %v, second half - %v, combined - %v\n", x, y, x+y)
	fmt.Println()

	ch := make(chan int, 2)
	ch <- 1
	ch <- 2
	//ch <- 3 fatal error: all goroutines are asleep - deadlock! - goroutine 1 [chan send]
	fmt.Println(<-ch)
	fmt.Println(<-ch)
	//fmt.Println(<-ch) fatal error: all goroutines are asleep - deadlock! - goroutine 1 [chan receive]

}

// function sum summarizes all values inside argument's slice and sends the summary to integer channel
func sum(s []int, c chan int) {
	sum := 0
	for _, v := range s {
		sum += v //calculating summary
	}
	c <- sum //sending summary to the int channel from argument
}
