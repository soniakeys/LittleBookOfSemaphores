package main

import (
	"log"
	"sync"
)

const nGR = 5

var wg sync.WaitGroup
var barrier = make(chan int)

func gr(grn int) {
	log.Println("gr", grn, "rendezvous (before barrier)")
	wg.Done()
	<-barrier
	log.Println("gr", grn, "critical point (after barrier)")
	wg.Done()
}

func main() {
	wg.Add(nGR)
	for i := 0; i < nGR; i++ {
		go gr(i)
	}
	wg.Wait()
	wg.Add(nGR)
	close(barrier)
	wg.Wait()
}
