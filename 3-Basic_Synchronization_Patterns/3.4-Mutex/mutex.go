package main

import (
	"fmt"
	"sync"
)

var count int

func inc(mutex *sync.Mutex, wg *sync.WaitGroup) {
	mutex.Lock()
	count++
	mutex.Unlock()
	wg.Done()
}

func main() {
	var mutex sync.Mutex
	var wg sync.WaitGroup
	wg.Add(2)
	go inc(&mutex, &wg)
	go inc(&mutex, &wg)
	wg.Wait()
	fmt.Println("count:", count)
}
