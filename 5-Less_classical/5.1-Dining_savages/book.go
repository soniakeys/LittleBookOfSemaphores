package main

import (
	"log"

	"github.com/soniakeys/LittleBookOfSemaphores/sem"
)

var (
	servings = 0
	mutex    = sem.NewChanSem(1, 1)
	emptyPot = sem.NewChanSem(0, 1)
	fullPot  = sem.NewChanSem(0, 1)
)

const M = 4
const nSavages = 6
const nPots = 3

func savage(n int) {
	for {
		mutex.Wait()
		if servings == 0 {
			emptyPot.Signal()
			fullPot.Wait()
			servings = M
		}
		servings -= 1
		getServingFromPot(n)
		mutex.Signal()
		eat(n)
	}
}

func getServingFromPot(n int) {
	log.Print("savage ", n, " gets serving from pot")
}

func eat(n int) {
	log.Print("savage ", n, " eats")
}

func main() {
	log.Print(nPots, " pot fillings:")
	for i := 1; i <= nSavages; i++ {
		go savage(i)
	}
	for j := 0; ; j++ {
		emptyPot.Wait()
		if j == nPots {
			return
		}
		putServingsInPot(M)
		fullPot.Signal()
	}
}

func putServingsInPot(m int) {
	log.Print("cook puts ", m, " servings in pot")
}
