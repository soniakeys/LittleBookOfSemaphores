package main

import (
	"log"
	"sync"

	"github.com/soniakeys/LittleBookOfSemaphores/sem"
)

var (
	n       = 5
	count   = 0
	mutex   = sem.NewChanSem(1, 1)
	barrier = sem.NewChanSem(0, 1)
)

var wg sync.WaitGroup

func gr(grn int, barrier sem.ChanSem) {
	log.Println("gr", grn, "rendezvous")
	mutex.Wait()
	count = count + 1
	mutex.Signal()
	if count == n {
		barrier.Signal()
	}
	barrier.Wait()
	barrier.Signal()
	log.Println("gr", grn, "critical point")
	wg.Done()
}

func main() {
	wg.Add(n)
	for i := 0; i < n; i++ {
		go gr(i, barrier)
	}
	wg.Wait()
}
