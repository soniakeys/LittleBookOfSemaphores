package main

import (
	"fmt"
	"sync"

	"github.com/soniakeys/LittleBookOfSemaphores/sem"
)

var count int

func inc(mutex sem.ChanSem) {
	mutex.Wait()
	count++
	mutex.Signal()
	wg.Done()
}

var wg sync.WaitGroup

func main() {
	mutex := sem.NewChanSem(1, 1)
	wg.Add(2)
	go inc(mutex)
	go inc(mutex)
	wg.Wait()
	fmt.Println("count:", count)
}
