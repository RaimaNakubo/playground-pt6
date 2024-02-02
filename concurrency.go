/*
Exercise: Equivalent Binary Trees
1. Implement the Walk function.

2. Test the Walk function.

The function tree.New(k) constructs a randomly-structured (but always sorted) binary tree holding the values k, 2k, 3k, ..., 10k.

Create a new channel ch and kick off the walker:

go Walk(tree.New(1), ch)
Then read and print 10 values from the channel. It should be the numbers 1, 2, 3, ..., 10.

3. Implement the Same function using Walk to determine whether t1 and t2 store the same values.

4. Test the Same function.

Same(tree.New(1), tree.New(1)) should return true, and Same(tree.New(1), tree.New(2)) should return false.

Doc for Tree package - https://pkg.go.dev/golang.org/x/tour/tree#Tree
*/

package main

import (
	"fmt"

	"golang.org/x/tour/tree"
)

func main() {

	//testing Walk() in another goroutine
	ch := make(chan int)
	go Walk(tree.New(1), ch)
	for i := range ch {
		fmt.Println(i)
	}
	fmt.Println()

	ch = make(chan int) //have to create a new channel bc Walk() closes it
	go Walk(tree.New(3), ch)
	for i := range ch {
		fmt.Println(i)
	}
	fmt.Println()

	//testing Same() with not necessarely identical trees, but with the same values for nodes total;
	//different arrangement, but same values
	k := 3                                // k is a multiplier for values of a tree
	res := Same(tree.New(k), tree.New(k)) //exp r true
	fmt.Println(res)

	//different arrangement, different values
	res = Same(tree.New(k), tree.New(1)) //exp r false
	fmt.Println(res)
}

// function Walk() walks the tree t sending all values from the tree to the channel ch and closes ch.
func Walk(t *tree.Tree, ch chan int) {
	defer close(ch) //before Walk() returns close the channel

	var nodeSearch func(t *tree.Tree) //defining a recursive function that searches for existing nodes in the tree

	//nodeSearch acts like a container for all the recursive calls
	nodeSearch = func(t *tree.Tree) {
		//if there is no node then search stops
		if t == nil {
			//fmt.Println("There is no node")
			return
		} /*else {
			fmt.Printf("There is a node: %v\n", t.Value)
		}

		if t.Left != nil {
			fmt.Printf("There is a node on the left: %v\n", t.Left.Value)
		} else {
			fmt.Println("There is no node on the left")
		}
		if t.Right != nil {
			fmt.Printf("There is a node on the right: %v\n", t.Right.Value)
		} else {
			fmt.Println("There is no node on the right")
		}*/

		//if there is a node then
		nodeSearch(t.Left) //search for node on the left
		ch <- t.Value      //send node's value to the channel
		//fmt.Printf("a node have been sent to the channel: %v\n", t.Value)
		nodeSearch(t.Right) //search for node on the right
	}

	nodeSearch(t) //calling recursive function to search for nodes in tree t

}

// function Same() determines whether the trees t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
	//channels ch1 and ch2 are recievers for data from binary trees
	ch1 := make(chan int)
	ch2 := make(chan int)

	//calling Walk() in separate goroutines which simultaneously sends data from trees to channels
	go Walk(t1, ch1)
	go Walk(t2, ch2)

	for v1 := range ch1 { //while first channel is open
		if v1 != <-ch2 { //compare all recieved values one-by-one (operator awaits value from ch2)
			// if ch2 is closed before ch1 closes then panic occurs; add channel's status comparison to handle this case if trees have different length
			return false //if any pair doesn't match then trees are not the same
		}
	}

	//if all recieved values from ch1 matches corresponding values recieved from ch2,
	//then trees are the same
	return true
}
