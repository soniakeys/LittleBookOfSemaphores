package main

import (
	"log"
	"sync"
)

const nCycles = 4
const nGR = 3

var wg sync.WaitGroup

func gr(grn int, workDesc string) {
	log.Println("gr", grn, workDesc)
	wg.Done()
}

func main() {
	wg.Add(nGR)
	for i := 0; i < nGR; i++ {
		go gr(i, "before-work")
	}
	wg.Wait()

	wg.Add(nGR)
	for i := 0; i < nGR; i++ {
		go gr(i, "after-work")
	}
	wg.Wait()
}
