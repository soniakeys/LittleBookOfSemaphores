package main

import (
	"log"
	"sync"
	"sync/atomic"

	"github.com/soniakeys/LittleBookOfSemaphores/sem"
)

var (
	mutex  sem.ChanSem
	items  sem.ChanSem
	buffer eventBuffer // not in hint, but this is what mutex protects
)

func producer(np int) {
	event := WaitForEvent()
	log.Println("producer", np, "produces event", event)
	mutex.Wait()
	buffer.add(event)
	mutex.Signal()
	items.Signal()
}

func consumer(nc int) {
	for {
		items.Wait()
		mutex.Wait()
		event := buffer.get()
		mutex.Signal()
		log.Println("  consumer", nc, "gets event", event)
		event.process()
	}
}

var last int64 // last event number generated
var wg sync.WaitGroup

type event64 int64

func (e event64) process() {
	log.Println("    processed: event", e)
	wg.Done()
}

func WaitForEvent() event64 {
	return event64(atomic.AddInt64(&last, 1))
}

type eventBuffer []event64

func (b *eventBuffer) add(e event64) {
	*b = append(*b, e)
}

func (b *eventBuffer) get() (e event64) {
	e = (*b)[0]
	*b = (*b)[1:]
	return
}

const nEvents = 6

func main() {
	mutex = sem.NewChanSem(1, 1)
	items = sem.NewChanSem(0, nEvents)
	go consumer(1)
	go consumer(2)
	wg.Add(nEvents)
	for i := 1; i <= nEvents; i++ {
		go producer(i)
	}
	wg.Wait()
}
