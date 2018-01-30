package main

import (
	"bytes"
	"log"
	"sync"
)

var (
	m  sync.RWMutex
	b  bytes.Buffer
	wg sync.WaitGroup
)

func writer(nw int) {
	m.Lock()
	log.Println("writer", nw, "writes")
	b.WriteString("w")
	m.Unlock()
	wg.Done()
}

func reader(nr int) {
	m.RLock()
	log.Println("reader", nr, "sees", b.Len(), "bytes")
	m.RUnlock()
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
