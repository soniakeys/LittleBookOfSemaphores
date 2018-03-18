package main

import (
	"log"
	"sync"
)

var l0 = make(chan int)
var l1 = make(chan int)
var u0 = make(chan int)
var u1 = make(chan int)
var wg sync.WaitGroup

func car() {
	for {
		load()
		for i := 0; i < C; i++ {
			<-l0
		}
		for i := 0; i < C; i++ {
			<-l1
		}
		run()
		unload()
		for i := 0; i < C; i++ {
			<-u0
		}
		for i := 0; i < C; i++ {
			<-u1
		}
	}
}

func load()   { log.Print("car ready to load") }
func run()    { log.Print("car runs") }
func unload() { log.Print("car ready to unload") }

func passenger() {
	l0 <- 1
	board()
	l1 <- 1
	u0 <- 1
	unboard()
	u1 <- 1
	wg.Done()
}

func board()   { log.Print("passenger boards") }
func unboard() { log.Print("passenger unboards") }

const (
	C           = 4
	nPassengers = 12
)

func main() {
	go car()
	wg.Add(nPassengers)
	for i := 0; i < nPassengers; i++ {
		go passenger()
	}
	wg.Wait()
}
