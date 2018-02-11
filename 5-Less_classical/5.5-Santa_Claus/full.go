package main

import (
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

const (
	nElves          = 8000 // race detector limit is 8192 concurrent goroutines
	meanWorkTime    = 1e4 * time.Millisecond
	meanHelpTime    = 3 * time.Millisecond
	vacationTime    = 20 * time.Millisecond
	meanReturnStall = 3 * time.Millisecond
)

var (
	reindeer = []string{"Dasher", "Dancer", "Prancer", "Vixen",
		"Comet", "Cupid", "Donner", "Blitzen", "Rudolph"}
	sleighReady = make(chan int)
	dashAway    = make(chan int)
	hitched     sync.WaitGroup

	nReindeer       int64
	first           = make(chan string)
	reindeerAllHere = make(chan string)

	nElvesNeedHelp      int64
	elvesNeedHelp       = make(chan elfNeedingHelp)
	thirdElfNeedingHelp = make(chan int, int(nElves)/3)
	elvesHelped         sync.WaitGroup
)

type elfNeedingHelp struct {
	n            int
	santaHelping chan int
}

func main() {
	go santa()
	for _, r := range reindeer {
		go reindeerFunc(r)
	}
	for i := 1; i <= nElves; i++ {
		go elfFunc(i)
	}
	<-dashAway
	log.Print("Christmas")
}

func santa() {
	helpElves := func() {
		log.Print("Santa lets in three elves")
		e1 := <-elvesNeedHelp
		e2 := <-elvesNeedHelp
		e3 := <-elvesNeedHelp
		elvesHelped.Add(3)
		e1.santaHelping <- 1
		e2.santaHelping <- 1
		e3.santaHelping <- 1
		log.Printf("Santa helps elves %d, %d, and %d", e1.n, e2.n, e3.n)
		elvesHelped.Wait()
		log.Print("Santa lets three elves out")
		atomic.AddInt64(&nElvesNeedHelp, -3)
	}
sleep:
	for {
		log.Print("Santa sleeping")
		select {
		case r := <-reindeerAllHere:
			log.Print("Santa awakened by ", r)
			break sleep
		default:
			select {
			case r := <-reindeerAllHere:
				log.Print("Santa awakened by ", r)
				break sleep
			case e := <-thirdElfNeedingHelp:
				log.Print("Santa awakened by elf ", e)
				helpElves()
			}
		}
		// Santa awake in this loop
		for {
			select {
			case r := <-reindeerAllHere:
				log.Print(r,
					" arrives to tell Santa reindeer are all here")
				break sleep
			default:
				select {
				case <-thirdElfNeedingHelp:
					helpElves()
				default:
					continue sleep
				}
			}
		}
	}
	log.Print("Santa prepares sleigh")
	hitched.Add(9)
	close(sleighReady)
	hitched.Wait()
	dashAway <- 1
}

func reindeerFunc(r string) {
	time.Sleep(vacationTime + time.Duration(rand.Intn(2*int(meanReturnStall))))
	log.Print(r, " back from tropics")
	switch atomic.AddInt64(&nReindeer, 1) {
	case 1:
		first <- r
		log.Print(r, " waits for others in warming hut")
	case 2:
		log.Print(r, " joins ", <-first, " in warming hut")
	default:
		log.Print(r, " joins others in warming hut")
	case 9:
		log.Print(r, " waking Santa")
		reindeerAllHere <- r
	}
	<-sleighReady
	log.Print(r, " gets hitched to sleigh")
	hitched.Done()
}

func elfFunc(e int) {
	for {
		time.Sleep(time.Duration(rand.Intn(int(meanWorkTime))))
		log.Print("elf ", e, " needs help")
		if atomic.AddInt64(&nElvesNeedHelp, 1)%3 == 0 {
			log.Print("elf ", e, " goes to wake Santa")
			thirdElfNeedingHelp <- e
		}
		need := elfNeedingHelp{e, make(chan int)}
		elvesNeedHelp <- need
		<-need.santaHelping
		log.Print("elf ", e, " getting help")
		time.Sleep(time.Duration(rand.Intn(int(meanHelpTime))))
		elvesHelped.Done()
	}
}
