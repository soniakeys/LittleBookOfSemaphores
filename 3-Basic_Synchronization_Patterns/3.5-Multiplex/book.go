package main

import (
	"fmt"
	"sync"

	"github.com/soniakeys/LittleBookOfSemaphores/sem"
)

var count int

func inc(m sem.ChanSem) {
	m.Wait()
	count++
	m.Signal()
	wg.Done()
}

var wg sync.WaitGroup

func main() {
	nAllowed := 1000
	nGR := 5000
	m := sem.NewChanSem(nAllowed, nAllowed)
	wg.Add(nGR)
	for i := 0; i < nGR; i++ {
		go inc(m)
	}
	wg.Wait()
	fmt.Println("count:", count)
}
