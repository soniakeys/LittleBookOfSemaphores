package main

import (
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

var (
	servingsInPot        = 0
	potsCooked           = 0
	headOfLine           sync.Mutex
	seeEmptyPotWakeCook  = make(chan int)
	seeFullPotCookSleeps = make(chan int)
)

const M = 4
const nSavages = 6
const nPots = 3

func savage(n int) {
	for {
		time.Sleep(time.Duration(rand.Intn(1e6)))
		headOfLine.Lock()
		if servingsInPot == 0 {
			if potsCooked == nPots {
				os.Exit(0)
			}
			log.Print("savage ", n, " finds pot empty, wakes cook")
			seeEmptyPotWakeCook <- 1
			<-seeFullPotCookSleeps
		}
		servingsInPot--
		headOfLine.Unlock()
		log.Print("savage ", n, " eats")
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	log.Print(nPots, " pot fillings:")
	for i := 1; i <= nSavages; i++ {
		go savage(i)
	}
	for j := 0; ; j++ {
		<-seeEmptyPotWakeCook
		potsCooked++
		servingsInPot = M
		log.Print("cook puts ", M, " servings in pot")
		seeFullPotCookSleeps <- 1
	}
}
