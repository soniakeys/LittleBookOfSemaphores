package main

import (
	"log"
	"sync"
)

const nCycles = 4
const nGR = 3

var cycleStart, cycleReset chan int
var wg sync.WaitGroup
var allDone = make(chan int)

func gr(grn int) {
	count := 0
	for {
		select {
		case <-allDone:
			log.Println("gr", grn, "counted to", count)
			wg.Done()
			return
		case <-cycleStart:
			log.Println("  gr", grn, "working")
			count++
			wg.Done()
			<-cycleReset
		}
	}
}

func main() {
	for i := 0; i < nGR; i++ {
		go gr(i)
	}
	cycleStart = make(chan int)
	for i := 0; i < nCycles; i++ {
		log.Println("cycle", i)
		wg.Add(nGR)
		cycleReset = make(chan int)
		close(cycleStart)

		wg.Wait()
		cycleStart = make(chan int)
		close(cycleReset)
	}
	wg.Add(nGR)
	close(allDone)
	wg.Wait()
}
