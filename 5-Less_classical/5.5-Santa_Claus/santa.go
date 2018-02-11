package main

import (
	"log"
	"math/rand"
	"sync/atomic"
	"time"
)

const nElves = 8000

func main() {
	for _, r := range []string{"Dasher", "Dancer", "Prancer", "Vixen",
		"Comet", "Cupid", "Donner", "Blitzen", "Rudolph"} {
		go reindeer(r)
	}
	for i := 1; i <= nElves; i++ {
		go elf(i)
	}
sleep:
	for {
		select {
		case <-reindeerAllHere:
			break sleep
		default:
			select {
			case <-reindeerAllHere:
				break sleep
			case <-threeElvesNeedHelp:
				log.Printf("santa helps elves %d, %d, and %d",
					<-elvesNeedHelp, <-elvesNeedHelp, <-elvesNeedHelp)
				atomic.AddInt64(&nElvesNeedHelp, -3)
			}
		}
	}
	log.Print("christmas")
}

var nReindeer int64
var reindeerAllHere = make(chan int)

func reindeer(r string) {
	time.Sleep(time.Duration(rand.Intn(1e7)))
	log.Print(r, " back from tropics")
	if atomic.AddInt64(&nReindeer, 1) == 9 {
		log.Print(r, " waking santa")
		reindeerAllHere <- 1
	}
}

var nElvesNeedHelp int64
var elvesNeedHelp = make(chan int)
var threeElvesNeedHelp = make(chan struct{}, int(nElves)/3)

func elf(e int) {
	for {
		time.Sleep(time.Duration(rand.Intn(5e9)))
		log.Print("elf ", e, " needs help")
		if atomic.AddInt64(&nElvesNeedHelp, 1)%3 == 0 {
			log.Print("elf ", e, " goes to wake santa")
			threeElvesNeedHelp <- struct{}{}
		}
		elvesNeedHelp <- e
	}
}
