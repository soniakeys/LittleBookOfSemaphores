package main

import (
	"log"
	"sync"

	"github.com/soniakeys/LittleBookOfSemaphores/sem"
)

var (
	n          = 5
	count      = 0
	turnstile  = sem.NewChanSem(0, n)
	turnstile2 = sem.NewChanSem(1, n)
	mutex      = sem.NewChanSem(1, 1)
)

var wg sync.WaitGroup

func gr(grn int) {
	log.Println("gr", grn, "rendezvous")

	mutex.Wait()
	count++
	if count == n {
		turnstile.SignalN(n) // unlock the first
	}
	mutex.Signal()
	turnstile.Wait() // first turnstile

	log.Println("gr", grn, "critical point")

	mutex.Wait()
	count--
	if count == 0 {
		turnstile2.SignalN(n) // unlock the second
	}
	mutex.Signal()
	turnstile2.Wait() // second turnstile
	wg.Done()
}

func main() {
	wg.Add(n)
	for i := 1; i <= n; i++ {
		go gr(i)
	}
	wg.Wait()
}
