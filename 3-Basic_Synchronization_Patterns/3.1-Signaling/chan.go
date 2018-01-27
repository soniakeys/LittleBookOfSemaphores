package main

import "log"

func main() {
	a1done := make(chan int)
	go func() { // goroutine "A"
		log.Print("statement a1")
		a1done <- 1
	}()
	// goroutine "B"
	<-a1done
	log.Print("statement b1")
}
