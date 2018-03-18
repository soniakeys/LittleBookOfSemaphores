package main

import (
	"log"
	"sync"

	"github.com/soniakeys/LittleBookOfSemaphores/sem"
)

var (
	mutex        = sem.NewChanSem(1, 1)
	mutex2       = sem.NewChanSem(1, 1)
	boarders     = 0
	unboarders   = 0
	boardQueue   = sem.NewChanSem(0, 1)
	unboardQueue = sem.NewChanSem(0, 1)
	allAboard    = sem.NewChanSem(0, 1)
	allAshore    = sem.NewChanSem(0, 1)
)

func car() {
	for {
		load()
		boardQueue.SignalN(C)
		allAboard.Wait()
		run()
		unload()
		unboardQueue.SignalN(C)
		allAshore.Wait()
	}
}

func load()   { log.Print("car ready to load") }
func run()    { log.Print("car runs") }
func unload() { log.Print("car ready to unload") }

func passenger() {
	boardQueue.Wait()
	board()
	mutex.Wait()
	boarders++
	if boarders == C {
		allAboard.Signal()
		boarders = 0
	}
	mutex.Signal()
	unboardQueue.Wait()
	unboard()
	mutex2.Wait()
	unboarders++
	if unboarders == C {
		allAshore.Signal()
		unboarders = 0
	}
	mutex2.Signal()
	wg.Done()
}

func board()   { log.Print("passenger boards") }
func unboard() { log.Print("passenger unboards") }

var wg sync.WaitGroup

const (
	C           = 4
	nPassengers = 12
)

func main() {
	go car()
	wg.Add(nPassengers)
	for i := 0; i < nPassengers; i++ {
		go passenger()
	}
	wg.Wait()
}
