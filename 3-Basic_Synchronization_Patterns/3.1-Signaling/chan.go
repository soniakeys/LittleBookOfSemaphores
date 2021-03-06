package main

import "fmt"

func main() {
	a1done := make(chan int)
	go func() { // goroutine "A"
		fmt.Println("statement a1")
		a1done <- 1
	}()
	// goroutine "B"
	<-a1done
	fmt.Println("statement b1")
}
