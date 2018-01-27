package main

import (
	"log"
	"sync"
)

const nCycles = 4
const nGR = 3

var wg sync.WaitGroup

func gr(cycle, grn int) {
	log.Println("cycle", cycle, "gr", grn, "working")
	wg.Done()
}

func main() {
	for i := 0; i < nCycles; i++ {
		wg.Add(nGR)
		for j := 0; j < nGR; j++ {
			go gr(i, j)
		}
		wg.Wait()
	}
}
