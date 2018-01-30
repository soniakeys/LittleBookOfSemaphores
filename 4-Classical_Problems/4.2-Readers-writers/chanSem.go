package main

import (
	"bytes"
	"log"
	"sync"

	"github.com/soniakeys/LittleBookOfSemaphores/sem"
)

var (
	readers   = 0
	mutex     = sem.NewChanSem(1, 1)
	roomEmpty = sem.NewChanSem(1, 1)
	b         bytes.Buffer
	wg        sync.WaitGroup
)

func writer(nw int) {
	roomEmpty.Wait()
	log.Println("writer", nw, "writes")
	b.WriteString("w")
	roomEmpty.Signal()
	wg.Done()
}

func reader(nr int) {
	mutex.Wait()
	readers++
	if readers == 1 {
		roomEmpty.Wait() // first in locks
	}
	mutex.Signal()
	log.Println("reader", nr, "sees", b.Len(), "bytes")
	mutex.Wait()
	readers--
	if readers == 0 {
		roomEmpty.Signal() // last out unlocks
	}
	mutex.Signal()
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
