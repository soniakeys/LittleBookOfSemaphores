package main

import (
	"fmt"
	"sync"
)

var (
	count int
	mutex sync.Mutex
)

func inc(wg *sync.WaitGroup) {
	mutex.Lock()
	count++
	mutex.Unlock()
	wg.Done()
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	go inc(&wg)
	go inc(&wg)
	wg.Wait()
	fmt.Println("count:", count)
}
