package main

import (
	"fmt"

	"github.com/soniakeys/LittleBookOfSemaphores/sem"
)

func main() {
	a1done := sem.NewChanSem(0, 1)
	go func() { // goroutine "A"
		fmt.Println("statement a1")
		a1done.Signal()
	}()
	// goroutine "B"
	a1done.Wait()
	fmt.Println("statement b1")
}
