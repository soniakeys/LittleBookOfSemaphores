package main

import (
	"log"
	"sync"

	"github.com/soniakeys/LittleBookOfSemaphores/sem"
)

var (
	mutex      = sem.NewChanSem(1, 1)
	oxygen     = 0
	hydrogen   = 0
	barrier    = sem.NewBarrier(3)
	oxyQueue   = sem.NewChanSem(0, 1)
	hydroQueue = sem.NewChanSem(0, 1)
)

func oxyFunc() {
	mutex.Wait()
	oxygen++
	log.Print("oxygen count to ", oxygen)
	if hydrogen >= 2 {
		hydroQueue.SignalN(2)
		hydrogen -= 2
		oxyQueue.Signal()
		oxygen--
	} else {
		mutex.Signal()
	}
	oxyQueue.Wait()
	bond("O")
	barrier.Wait()
	mutex.Signal()
	wg.Done()
}

func hydroFunc() {
	mutex.Wait()
	hydrogen++
	log.Print("hydrogen count to ", hydrogen)
	if hydrogen >= 2 && oxygen >= 1 {
		hydroQueue.SignalN(2)
		hydrogen -= 2
		oxyQueue.Signal()
		oxygen--
	} else {
		mutex.Signal()
	}
	hydroQueue.Wait()
	bond("H")
	barrier.Wait()
	wg.Done()
}

func bond(e string) {
	log.Print(e, " bonds")
}

const (
	nO = 3
	nH = 2 * nO
)

var wg sync.WaitGroup

func main() {
	wg.Add(nH + nO)
	for i := 0; i < nH; i++ {
		go hydroFunc()
	}
	for i := 0; i < nO; i++ {
		go oxyFunc()
	}
	wg.Wait()
}
