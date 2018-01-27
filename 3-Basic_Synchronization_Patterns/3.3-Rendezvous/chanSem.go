package main

import (
	"log"
	"sync"

	"github.com/soniakeys/LittleBookOfSemaphores/sem"
)

func gr(name string, IArrived, OtherArrived sem.ChanSem) {
	log.Print("statement ", name, "1")
	IArrived.Signal()
	OtherArrived.Wait()
	log.Print("statement ", name, "2")
	wg.Done()
}

var wg sync.WaitGroup

func main() {
	aArrived := sem.NewChanSem(0, 1)
	bArrived := sem.NewChanSem(0, 1)
	wg.Add(2)
	go gr("a", aArrived, bArrived)
	go gr("b", bArrived, aArrived)
	wg.Wait()
}
