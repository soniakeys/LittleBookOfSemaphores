package main

import (
	"log"
	"sync"
	"sync/atomic"

	"github.com/soniakeys/LittleBookOfSemaphores/sem"
)

func producer(np int) {
	event := WaitForEvent()
	log.Println("producer", np, "produces event", event)
	spaces.Wait()
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
		spaces.Signal()
		log.Println("  consumer", nc, "gets event", event)
		event.process()
	}
}

var (
	last   int64 // last event number
	mutex  sem.ChanSem
	items  sem.ChanSem
	spaces sem.ChanSem
	buffer *eventBuffer
	wg     sync.WaitGroup
)

type event64 int64

func (e event64) process() {
	log.Println("    processed: event", e)
	wg.Done()
}

func WaitForEvent() event64 {
	return event64(atomic.AddInt64(&last, 1))
}

type eventBuffer struct {
	b          []event64
	start, end int
}

func newEventBuffer(size int) *eventBuffer {
	return &eventBuffer{b: make([]event64, size)}
}

func (b *eventBuffer) add(e event64) {
	b.b[b.end] = e
	b.end++
	if b.end == len(b.b) {
		b.end = 0
	}
}

func (b *eventBuffer) get() (e event64) {
	e = b.b[b.start]
	b.start++
	if b.start == len(b.b) {
		b.start = 0
	}
	return
}

const nEvents = 6
const size = 3

func main() {
	mutex = sem.NewChanSem(1, 1)
	items = sem.NewChanSem(0, nEvents)
	spaces = sem.NewChanSem(size, size)
	buffer = newEventBuffer(size)
	go consumer(1)
	go consumer(2)
	wg.Add(nEvents)
	for i := 1; i <= nEvents; i++ {
		go producer(i)
	}
	wg.Wait()
}
