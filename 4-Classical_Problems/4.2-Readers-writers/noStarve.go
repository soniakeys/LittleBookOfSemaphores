package main

import (
	"bytes"
	"log"
	"sync"

	"github.com/soniakeys/LittleBookOfSemaphores/sem"
)

var (
	readSwitch = sem.NewLightswitch()
	roomEmpty  = sem.NewChanSem(1, 1)
	turnstile  = sem.NewChanSem(1, 1)
	b          bytes.Buffer // protected by roomEmpty
)

var wg sync.WaitGroup

func writer(nw int) {
	turnstile.Wait()
	roomEmpty.Wait()
	log.Println("writer", nw, "writes")
	b.WriteString("w")
	turnstile.Signal()
	roomEmpty.Signal()
	wg.Done()
}

func reader(nr int) {
	turnstile.Wait()
	turnstile.Signal()
	readSwitch.Lock(roomEmpty)
	log.Println("reader", nr, "sees", b.Len(), "bytes")
	readSwitch.Unlock(roomEmpty)
	wg.Done()
}

const nw = 6
const nr = 6

func main() {
	wg.Add(nw + nr)
	for i := 1; i <= 6; i++ {
		go writer(i)
	}
	for i := 1; i <= 6; i++ {
		go reader(i)
	}
	wg.Wait()
}
