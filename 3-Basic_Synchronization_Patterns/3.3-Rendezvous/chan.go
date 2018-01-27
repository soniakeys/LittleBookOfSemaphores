package main

import "log"

func main() {
	ch := make(chan int)
	go func() { // "Thread" A
		log.Print("statement a1")
		ch <- 1
		log.Print("statement a2")
	}()
	// "Thread" B
	log.Print("statement b1")
	<-ch
	log.Print("statement b2")
}
