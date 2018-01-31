package main

import (
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/soniakeys/LittleBookOfSemaphores/sem"
)

var (
	fork    [5]sem.ChanSem
	rightie [5]bool
)

func init() {
	for i := range fork {
		fork[i] = sem.NewChanSem(1, 1)
	}
	rightie = [5]bool{true, true, true, true, false}
}

var wg sync.WaitGroup

const nBites = 4

func ph(n int) {
	log.Println("philosopher", n, "sits to eat")
	for b := 1; b <= nBites; b++ {
		think(n)
		get_forks(n)
		eat(n, b)
		put_forks(n)
	}
	log.Println("philosopher", n, "full, leaves table")
	wg.Done()
}

func think(n int) {
	log.Println("philosopher", n, "thinking")
	time.Sleep(time.Duration(rand.Intn(1e8)))
}

func eat(n, b int) {
	log.Printf("philosopher %d eats bite #%d", n, b)
	time.Sleep(time.Duration(rand.Intn(1e8)))
}

func get_forks(i int) {
	log.Println("philosopher", i, "wants to eat, looks for forks")
	if rightie[i] {
		fork[right(i)].Wait()
		log.Println("philosopher", i, "has right fork")
		fork[left(i)].Wait()
		log.Println("philosopher", i, "has left fork")
	} else {
		fork[left(i)].Wait()
		log.Println("philosopher", i, "has left fork")
		fork[right(i)].Wait()
		log.Println("philosopher", i, "has right fork")
	}
}

func put_forks(i int) {
	fork[right(i)].Signal()
	fork[left(i)].Signal()
	log.Println("philosopher", i, "satisfied for the moment, returns forks")
}

func right(i int) int { return i }
func left(i int) int  { return (i + 1) % 5 }

func main() {
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go ph(i)
	}
	wg.Wait()
}
