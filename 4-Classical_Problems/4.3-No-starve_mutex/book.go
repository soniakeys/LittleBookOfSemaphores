package main

import (
	"fmt"
	"sync"

	"github.com/soniakeys/LittleBookOfSemaphores/sem"
)

var (
	room1, room2 = 0, 0
	mutex        = sem.NewChanSem(1, 1)
	t1           = sem.NewChanSem(1, 1)
	t2           = sem.NewChanSem(0, 1)
)

var wg sync.WaitGroup

func gr(grn int) {
	mutex.Wait()
	room1++
	mutex.Signal()

	t1.Wait()
	room2++
	mutex.Wait()
	room1--
	if room1 == 0 {
		mutex.Signal()
		t2.Signal()
	} else {
		mutex.Signal()
		t1.Signal()
	}

	t2.Wait()
	room2--
	fmt.Println("gr", grn, "runs")
	if room2 == 0 {
		t1.Signal()
	} else {
		t2.Signal()
	}
	wg.Done()
}

func main() {
	const ngr = 10
	wg.Add(ngr)
	for i := 1; i <= ngr; i++ {
		go gr(i)
	}
	wg.Wait()
}
