package main

import (
	"bytes"
	"log"
	"sync"

	"github.com/soniakeys/LittleBookOfSemaphores/sem"
)

var (
	readSwitch  = sem.NewLightswitch()
	writeSwitch = sem.NewLightswitch()
	noReaders   = sem.NewChanSem(1, 1)
	noWriters   = sem.NewChanSem(1, 1)
	b           bytes.Buffer // protected by noWriters
)

var wg sync.WaitGroup

func writer(nw int) {
	writeSwitch.Lock(noReaders)
	noWriters.Wait()
	log.Println("writer", nw, "writes")
	b.WriteString("w")
	noWriters.Signal()
	writeSwitch.Unlock(noReaders)
	wg.Done()
}

func reader(nr int) {
	noReaders.Wait()
	readSwitch.Lock(noWriters)
	noReaders.Signal()
	log.Println("reader", nr, "sees", b.Len(), "bytes")
	readSwitch.Unlock(noWriters)
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
