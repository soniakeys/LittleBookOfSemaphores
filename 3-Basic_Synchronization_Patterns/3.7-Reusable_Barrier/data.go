package main

import (
	"log"
	"sync"
)

const nCycles = 4
const nGR = 3

var wg sync.WaitGroup

func gr(cycle, grn int, count *int) {
	log.Println("cycle", cycle, "gr", grn, "working")
	*count++
	wg.Done()
}

func main() {
	count := make([]int, nGR)
	for i := 0; i < nCycles; i++ {
		wg.Add(nGR)
		for j := 0; j < nGR; j++ {
			go gr(i, j, &count[j])
		}
		wg.Wait()
	}
	for j := 0; j < nGR; j++ {
		log.Printf("count[%d] = %d", j, count[j])
	}
}
