package main

import (
	"log"

	"github.com/soniakeys/LittleBookOfSemaphores/sem"
)

func main() {
	a1done := sem.NewCountSem(0)
	go func() { // goroutine "A"
		log.Print("statement a1")
		a1done.Signal()
	}()
	// goroutine "B"
	a1done.Wait()
	log.Print("statement b1")
}
