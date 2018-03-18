package main

import (
	"log"
	"sync"

	"github.com/soniakeys/LittleBookOfSemaphores/sem"
)

const m = 3

var (
	mutex         = sem.NewChanSem(1, 1)
	mutex2        = sem.NewChanSem(1, 1)
	boarders      = 0
	unboarders    = 0
	boardQueue    = sem.NewChanSem(0, 1)
	unboardQueue  = sem.NewChanSem(0, 1)
	allAboard     = sem.NewChanSem(0, 1)
	allAshore     = sem.NewChanSem(0, 1)
	loadingArea   [m]sem.ChanSem
	unloadingArea [m]sem.ChanSem
)

func init() {
	for i := 0; i < m; i++ {
		loadingArea[i] = sem.NewChanSem(0, 1)
		unloadingArea[i] = sem.NewChanSem(0, 1)
	}
	loadingArea[0].Signal()
	unloadingArea[0].Signal()
}

func next(i int) int { return (i + 1) % m }

func car(i int) {
	for {
		loadingArea[i].Wait()
		load(i)
		boardQueue.SignalN(C)
		allAboard.Wait()
		loadingArea[next(i)].Signal()
		run(i)
		unloadingArea[i].Wait()
		unload(i)
		unboardQueue.SignalN(C)
		allAshore.Wait()
		unloadingArea[next(i)].Signal()
	}
}

func load(i int)   { log.Print("car ", i, " ready to load") }
func run(i int)    { log.Print("car ", i, " runs") }
func unload(i int) { log.Print("car ", i, " ready to unload") }

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
	nPassengers = 16
)

func main() {
	for i := 0; i < m; i++ {
		go car(i)
	}
	wg.Add(nPassengers)
	for i := 0; i < nPassengers; i++ {
		go passenger()
	}
	wg.Wait()
}
