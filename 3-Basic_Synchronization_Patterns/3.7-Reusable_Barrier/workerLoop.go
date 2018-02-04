package main

import (
	"log"
	"sync"
)

const nCycles = 4
const nGR = 3

var cycleStart, cycleReset chan int
var wg sync.WaitGroup
var quit = make(chan int)

func gr(grn int) {
	count := 0
	for {
		select {
		case <-cycleStart:
			log.Println("  gr", grn, "working")
			count++
			wg.Done()

			<-cycleReset
			wg.Done()
		case <-quit:
			log.Println("gr", grn, "counted to", count)
			wg.Done()
			return
		}
	}
}

func main() {
	cycleStart = make(chan int)
	for i := 0; i < nGR; i++ {
		go gr(i)
	}
	for i := 0; i < nCycles; i++ {
		log.Println("cycle", i)
		wg.Add(nGR)
		cycleReset = make(chan int)
		close(cycleStart)
		wg.Wait()

		wg.Add(nGR)
		cycleStart = make(chan int)
		close(cycleReset)
		wg.Wait()
	}
	wg.Add(nGR)
	close(quit)
	wg.Wait()
}
