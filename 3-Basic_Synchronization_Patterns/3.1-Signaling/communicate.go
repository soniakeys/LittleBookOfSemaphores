package main

import "fmt"

func main() {
	a1msg := make(chan string)
	go func() { // goroutine "A"
		a1msg <- "statement a1"
	}()
	// goroutine "B"
	fmt.Println(<-a1msg)
	fmt.Println("statement b1")
}
