package main

import (
	"log"
	"sync"
	"sync/atomic"

	"github.com/soniakeys/LittleBookOfSemaphores/sem"
)

var ( // Agent semaphores
	agentSem = sem.NewChanSem(1, 1)
	tobacco  = sem.NewChanSem(0, 1)
	paper    = sem.NewChanSem(0, 1)
	match    = sem.NewChanSem(0, 1)
)

var ( // Smoker semaphores
	isTobacco, isPaper, isMatch = false, false, false

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
			agentSem.Wait()
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
			agentSem.Wait()
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
			agentSem.Wait()
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
			if isPaper {
				isPaper = false
				matchSem.Signal()
			} else if isMatch {
				isMatch = false
				paperSem.Signal()
			} else {
				isTobacco = true
			}
			mutex.Signal()
		}
	}()
	go func() { // Pusher B
		for {
			paper.Wait()
			mutex.Wait()
			if isTobacco {
				isTobacco = false
				matchSem.Signal()
			} else if isMatch {
				isMatch = false
				tobaccoSem.Signal()
			} else {
				isPaper = true
			}
			mutex.Signal()
		}
	}()
	go func() { // Pusher C
		for {
			match.Wait()
			mutex.Wait()
			if isPaper {
				isPaper = false
				tobaccoSem.Signal()
			} else if isTobacco {
				isTobacco = false
				paperSem.Signal()
			} else {
				isMatch = true
			}
			mutex.Signal()
		}
	}()
	go func() { // Smoker with tobacco
		for {
			tobaccoSem.Wait()
			makeCigarette("tobacco")
			agentSem.Signal()
			smoke("tobacco")
		}
	}()
	go func() { // Smoker with paper
		for {
			paperSem.Wait()
			makeCigarette("paper")
			agentSem.Signal()
			smoke("paper")
		}
	}()
	go func() { // Smoker with matches
		for {
			matchSem.Wait()
			makeCigarette("matches")
			agentSem.Signal()
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
