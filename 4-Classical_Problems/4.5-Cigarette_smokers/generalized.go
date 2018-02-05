package main

import (
	"log"
	"sync"
	"sync/atomic"

	"github.com/soniakeys/LittleBookOfSemaphores/sem"
)

var ( // Agent semaphores
	tobacco = sem.NewChanSem(0, 1)
	paper   = sem.NewChanSem(0, 1)
	match   = sem.NewChanSem(0, 1)
)

var ( // Smoker semaphores
	numTobacco, numPaper, numMatch = 0, 0, 0

	tobaccoSem = sem.NewChanSem(0, 1)
	paperSem   = sem.NewChanSem(0, 1)
	matchSem   = sem.NewChanSem(0, 1)
	mutex      = sem.NewChanSem(1, 1)
)

var provided int64
var wg sync.WaitGroup

const rounds = 10

func main() {
	wg.Add(rounds)
	go func() { // Agent A code
		for {
			if atomic.AddInt64(&provided, 1) > rounds {
				return
			}
			log.Println("agent provides tobacco and paper")
			tobacco.Signal()
			paper.Signal()
		}
	}()
	go func() { // Agent B code
		for {
			if atomic.AddInt64(&provided, 1) > rounds {
				return
			}
			log.Println("agent provides paper and a match")
			paper.Signal()
			match.Signal()
		}
	}()
	go func() { // Agent C code
		for {
			if atomic.AddInt64(&provided, 1) > rounds {
				return
			}
			log.Println("agent provides tobacco and a match")
			tobacco.Signal()
			match.Signal()
		}
	}()
	go func() { // Pusher A
		for {
			tobacco.Wait()
			mutex.Wait()
			if numPaper > 0 {
				numPaper--
				matchSem.Signal()
			} else if numMatch > 0 {
				numMatch--
				paperSem.Signal()
			} else {
				numTobacco++
			}
			mutex.Signal()
		}
	}()
	go func() { // Pusher B
		for {
			paper.Wait()
			mutex.Wait()
			if numTobacco > 0 {
				numTobacco--
				matchSem.Signal()
			} else if numMatch > 0 {
				numMatch--
				tobaccoSem.Signal()
			} else {
				numPaper++
			}
			mutex.Signal()
		}
	}()
	go func() { // Pusher C
		for {
			match.Wait()
			mutex.Wait()
			if numPaper > 0 {
				numPaper--
				tobaccoSem.Signal()
			} else if numTobacco > 0 {
				numTobacco--
				paperSem.Signal()
			} else {
				numMatch++
			}
			mutex.Signal()
		}
	}()
	go func() { // Smoker with tobacco
		for {
			tobaccoSem.Wait()
			makeCigarette("tobacco")
			smoke("tobacco")
		}
	}()
	go func() { // Smoker with paper
		for {
			paperSem.Wait()
			makeCigarette("paper")
			smoke("paper")
		}
	}()
	go func() { // Smoker with matches
		for {
			matchSem.Wait()
			makeCigarette("matches")
			smoke("matches")
		}
	}()
	wg.Wait()
}

func makeCigarette(s string) {
	log.Println("smoker with", s, "makes cigarette")
}

func smoke(s string) {
	log.Println("smoker with", s, "smokes")
	wg.Done()
}
