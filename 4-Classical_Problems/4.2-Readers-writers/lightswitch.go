package main

import (
	"bytes"
	"log"
	"sync"

	"github.com/soniakeys/LittleBookOfSemaphores/sem"
)

var (
	readLightswitch = sem.NewLightswitch()
	roomEmpty       = sem.NewChanSem(1, 1)
	b               bytes.Buffer
	wg              sync.WaitGroup
)

func writer(nw int) {
	roomEmpty.Wait()
	log.Println("writer", nw, "writes")
	b.WriteString("w")
	roomEmpty.Signal()
	wg.Done()
}

func reader(nr int) {
	readLightswitch.Lock(roomEmpty)
	log.Println("reader", nr, "sees", b.Len(), "bytes")
	readLightswitch.Unlock(roomEmpty)
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
