package main

import (
	"log"
	"sync"

	"github.com/soniakeys/LittleBookOfSemaphores/sem"
)

var (
	n          = 5
	count      = 0
	turnstile  = sem.NewChanSem(0, 1)
	turnstile2 = sem.NewChanSem(1, 1)
	mutex      = sem.NewChanSem(1, 1)
)

var wg sync.WaitGroup

func gr(grn int) {
	log.Println("gr", grn, "rendezvous")
	mutex.Wait()
	count++
	if count == n {
		turnstile2.Wait()  // lock the second
		turnstile.Signal() // unlock the first
	}
	mutex.Signal()
	turnstile.Wait() // first turnstile
	turnstile.Signal()
	log.Println("gr", grn, "critical point")
	mutex.Wait()
	count--
	if count == 0 {
		turnstile.Wait()    // lock the first
		turnstile2.Signal() // unlock the second
	}
	mutex.Signal()
	turnstile2.Wait() // second turnstile
	turnstile2.Signal()
	wg.Done()
}

func main() {
	wg.Add(n)
	for i := 1; i <= n; i++ {
		go gr(i)
	}
	wg.Wait()
}
