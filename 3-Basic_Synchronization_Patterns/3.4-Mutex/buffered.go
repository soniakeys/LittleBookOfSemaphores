package main

import (
	"fmt"
	"sync"
)

var count = make(chan int, 1)

func inc(wg *sync.WaitGroup) {
	count <- <-count + 1
	wg.Done()
}

func main() {
	count <- 0
	var wg sync.WaitGroup
	wg.Add(2)
	go inc(&wg)
	go inc(&wg)
	wg.Wait()
	fmt.Println("count:", <-count)
}
