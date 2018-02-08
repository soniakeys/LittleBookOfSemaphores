package main

import (
	"log"
	"sync"
	"sync/atomic"
	"time"
)

const M = 4
const nSavages = 6
const nPots = 3

var pot int64
var wake = make(chan chan int)
var wg sync.WaitGroup

func savage(n int) {
	watch := make(chan int)
	var s int64
	for {
		// log.Print("savage ", n, " out hunting/gathering")
		// time.Sleep(time.Duration(rand.Intn(1e7)))
		log.Print("savage ", n, " hungry") //, returns to camp")
		for {
			if s = atomic.AddInt64(&pot, -1); s >= 0 {
				break
			}
			log.Print("savage ", n, " finds pot empty, yells for cook")
			wake <- watch
			log.Print("savage ", n, " waiting to see full pot")
			<-watch
			log.Print("savage ", n, " sees servings in pot")
		}
		log.Print("savage ", n, " eats (leaving ", s, " servings in pot)")
		// time.Sleep(1e6)
		wg.Done()
	}
}

func main() {
	wg.Add(nPots * M)
	for i := 1; i <= nSavages; i++ {
		go savage(i)
	}
	awake := false
	waiting := make([]chan int, 0, nSavages)
	var stew <-chan time.Time
	for pots := 0; ; {
		select {
		case s := <-wake:
			// simulation code: "register" savage watching
			waiting = append(waiting, s)
			// cook wakes if not already awake
			if !awake {
				awake = true
				// cook must see pot empty
				if atomic.LoadInt64(&pot) > 0 {
					log.Println("cook grumbles, goes back to sleep")
					awake = false
					// simulation code: savage must now notice same non-empty
					// pot that cook sees, and stop waiting.
					s <- 1
					waiting = waiting[:len(waiting)-1]
				} else {
					log.Println("cook awake, starts cooking")
					stew = time.After(1e5)
				}
			}
		case <-stew:
			pots++
			log.Printf("cook puts %d servings in pot (%d pots cooked)", M, pots)
			atomic.StoreInt64(&pot, M)

			// simulation code: pot filling event observed by waiting savages
			for _, s := range waiting {
				s <- 1
			}
			waiting = waiting[:0]
			if pots == nPots {
				log.Println("cook leaves")
				wg.Wait() // simulation waits for savages to finish eating
				log.Println("simulation ends")
				return
			}
			log.Println("cook sleeps")
			awake = false
		}
	}
}
