package main

import (
	"fmt"
	"sync"
	"time"
)

// SafeCounter is used to count some strings and safe to use concurrently
type SafeCounter struct {
	mu sync.Mutex
	v  map[string]int
}

func main() {
	//c is an instance of our SafeCounter initialized with unlocked mutex(zero value) for mutex-lock and an empty map for value
	c := SafeCounter{v: make(map[string]int)}

	for i := 0; i < 1000; i++ { //for a thousand times
		go c.Inc("somekey") //incrementing counter for given key via separate goroutines
	} //so a thousand sub-goroutines will try to access value in counter

	//instantly after launching all sub-goroutines checking for a value in counter
	fmt.Println(c.Value("somekey")) //counter != 1000 bc some sub-goroutines cannot bypass locked mutex to increment counter in that short period of time

	//pausing main goroutine to wait for all sub-goroutines to complete their calls
	time.Sleep(time.Second)
	fmt.Println(c.Value("somekey")) //counter == 1000
}

// method Inc() safely per-goroutine increments counter from reciever for the key given in argument
func (c *SafeCounter) Inc(key string) {
	//locking mutex for map c.v so only one goroutine can access the c.v map at a time
	c.mu.Lock()

	c.v[key]++    //incrementing counter for given key
	c.mu.Unlock() //unlocking mutex
}

// method Value() safely per-goroutine returns the current value of the counter from reciever for a given key from argument
func (c *SafeCounter) Value(key string) int {
	c.mu.Lock()         //locking mutex for c.v map
	defer c.mu.Unlock() //mutex will unlock automaticly once this method returns

	return c.v[key] //returning counter for a given key
}
