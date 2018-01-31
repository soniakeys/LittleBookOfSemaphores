package main

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

const bites = 3

var dining sync.WaitGroup

func philosopher(ph int, dominantHand, otherHand *sync.Mutex) {
	log.Println("philospher", ph, "seated")
	rSleep := func() { time.Sleep(time.Duration(rand.Intn(1e8))) }
	for b := 1; b <= bites; b++ {
		log.Println("philospher", ph, "hungry")
		dominantHand.Lock() // pick up forks
		otherHand.Lock()
		log.Println("philospher", ph, "taking bite", b)
		rSleep()
		dominantHand.Unlock() // put down forks
		otherHand.Unlock()
		log.Println("philospher", ph, "thinking")
		rSleep()
	}
	log.Println("philospher", ph, "satisfied")
	dining.Done()
	log.Println("philospher", ph, "left the table")
}

func main() {
	log.Println("table empty")
	dining.Add(5)
	fork0 := &sync.Mutex{}
	forkLeft := fork0
	for i := 1; i < 5; i++ {
		forkRight := &sync.Mutex{}
		go philosopher(i, forkLeft, forkRight)
		forkLeft = forkRight
	}
	go philosopher(0, fork0, forkLeft)
	dining.Wait() // wait for philosphers to finish
	log.Println("table empty")
}
