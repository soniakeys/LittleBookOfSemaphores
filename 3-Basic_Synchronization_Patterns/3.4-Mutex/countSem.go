package main

import (
	"fmt"
	"sync"

	"github.com/soniakeys/LittleBookOfSemaphores/sem"
)

var count int

func inc(mutex *sem.CountSem) {
	mutex.Wait()
	count++
	mutex.Signal()
	wg.Done()
}

var wg sync.WaitGroup

func main() {
	mutex := sem.NewCountSem(1)
	wg.Add(2)
	go inc(mutex)
	go inc(mutex)
	wg.Wait()
	fmt.Println("count:", count)
}
