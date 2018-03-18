package main

import (
	"log"
	"sync"

	"github.com/soniakeys/LittleBookOfSemaphores/sem"
)

var (
	barrier     = sem.NewBarrier(4)
	mutex       = sem.NewChanSem(1, 1)
	hackers     = 0
	serfs       = 0
	hackerQueue = sem.NewChanSem(0, 1)
	serfQueue   = sem.NewChanSem(0, 1)
)

func hacker() {
	isCaptain := false
	mutex.Wait()
	hackers++
	if hackers == 4 {
		hackerQueue.SignalN(4)
		hackers = 0
		isCaptain = true
	} else if hackers == 2 && serfs >= 2 {
		hackerQueue.SignalN(2)
		serfQueue.SignalN(2)
		serfs -= 2
		hackers = 0
		isCaptain = true
	} else {
		mutex.Signal() // captain keeps the mutex
	}
	hackerQueue.Wait()
	board("hacker")
	barrier.Wait()
	if isCaptain {
		rowBoat("hacker")
		mutex.Signal() // captain releases the mutex
	}
}

func serf() {
	isCaptain := false
	mutex.Wait()
	serfs++
	if serfs == 4 {
		serfQueue.SignalN(4)
		serfs = 0
		isCaptain = true
	} else if serfs == 2 && hackers >= 2 {
		serfQueue.SignalN(2)
		hackerQueue.SignalN(2)
		hackers -= 2
		serfs = 0
		isCaptain = true
	} else {
		mutex.Signal() // captain keeps the mutex
	}
	serfQueue.Wait()
	board("serf")
	barrier.Wait()
	if isCaptain {
		rowBoat("serf")
		mutex.Signal() // captain releases the mutex
	}
}

func board(p string) {
	log.Print(p, " boards")
}

func rowBoat(p string) {
	log.Print(p, " rows")
	wg.Done()
}

const (
	nHackers = 6
	nSerfs   = 6
)

var wg sync.WaitGroup

func main() {
	wg.Add((nHackers + nSerfs) / 4)
	for i := 0; i < nHackers; i++ {
		go hacker()
	}
	for i := 0; i < nSerfs; i++ {
		go serf()
	}
	wg.Wait()
}
