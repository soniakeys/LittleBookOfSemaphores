package main

import (
	"log"
	"sync"
	"sync/atomic"

	"github.com/soniakeys/LittleBookOfSemaphores/sem"
)

const nGR = 5

var count int64
var wg sync.WaitGroup

func gr(grn int, barrier sem.ChanSem) {
	log.Println("gr", grn, "rendezvous")
	if c := atomic.AddInt64(&count, 1); c == nGR {
		barrier.Signal()
	}
	barrier.Wait()
	barrier.Signal()
	log.Println("gr", grn, "critical point")
	wg.Done()
}

func main() {
	barrier := sem.NewChanSem(0, 1)
	wg.Add(nGR)
	for i := 0; i < nGR; i++ {
		go gr(i, barrier)
	}
	wg.Wait()
}
