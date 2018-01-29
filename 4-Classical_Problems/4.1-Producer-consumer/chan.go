package main

import (
	"log"
	"sync"
	"sync/atomic"
)

func producer(np int) {
	event := WaitForEvent()
	log.Println("producer", np, "produces event", event)
	queue <- event
}

func consumer(nc int) {
	for {
		event := <-queue
		log.Println("  consumer", nc, "gets event", event)
		event.process()
	}
}

var (
	last  int64 // last event number
	queue = make(chan event64)
	wg    sync.WaitGroup
)

type event64 int64

func (e event64) process() {
	log.Println("    processed: event", e)
	wg.Done()
}

func WaitForEvent() event64 {
	return event64(atomic.AddInt64(&last, 1))
}

const nEvents = 6

func main() {
	go consumer(1)
	go consumer(2)
	wg.Add(nEvents)
	for i := 1; i <= nEvents; i++ {
		go producer(i)
	}
	wg.Wait()
}
