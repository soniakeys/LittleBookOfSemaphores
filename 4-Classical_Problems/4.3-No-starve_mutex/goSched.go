package main

import (
	"log"
	"sync"
)

var wg sync.WaitGroup

func gr(grn int) {
	log.Println("gr", grn, "runs")
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
