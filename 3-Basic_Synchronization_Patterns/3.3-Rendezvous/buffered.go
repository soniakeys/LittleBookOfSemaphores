package main

import (
	"log"
	"sync"
)

func gr(name string, IArrived, OtherArrived chan int) {
	log.Print("statement ", name, "1")
	IArrived <- 1
	<-OtherArrived
	log.Print("statement ", name, "2")
	wg.Done()
}

var wg sync.WaitGroup

func main() {
	aArrived := make(chan int, 1)
	bArrived := make(chan int, 1)
	wg.Add(2)
	go gr("a", aArrived, bArrived)
	go gr("b", bArrived, aArrived)
	wg.Wait()
}
