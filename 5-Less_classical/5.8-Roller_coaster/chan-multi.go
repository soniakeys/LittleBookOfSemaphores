package main

import (
	"log"
	"sync"
)

var l0 = make(chan int)
var l1 = make(chan int)
var u0 = make(chan int)
var u1 = make(chan int)

const C = 4
const m = 3

var track = make(chan int, m)
var platform = make(chan int, m)

func init() {
	for i := 0; i < m; i++ {
		platform <- i
	}
}

func depart() {
	for {
		car := <-platform
		load(car)
		for i := 0; i < C; i++ {
			l0 <- car
		}
		for i := 0; i < C; i++ {
			<-l1
		}
		run(car)
		track <- car
	}
}

func arrive() {
	for {
		car := <-track
		unload(car)
		for i := 0; i < C; i++ {
			<-u0
		}
		for i := 0; i < C; i++ {
			<-u1
		}
		platform <- car
	}
}

func load(i int)   { log.Print("car ", i, " ready to load") }
func run(i int)    { log.Print("car ", i, " runs") }
func unload(i int) { log.Print("car ", i, " ready to unload") }

func passenger() {
	car := <-l0
	board(car)
	l1 <- 1
	u0 <- 1
	unboard(car)
	u1 <- 1
	wg.Done()
}

func board(car int)   { log.Print("passenger boards car ", car) }
func unboard(car int) { log.Print("passenger unboards car ", car) }

const nPassengers = 12

var wg sync.WaitGroup

func main() {
	go depart()
	go arrive()
	wg.Add(nPassengers)
	for i := 0; i < nPassengers; i++ {
		go passenger()
	}
	wg.Wait()
}
