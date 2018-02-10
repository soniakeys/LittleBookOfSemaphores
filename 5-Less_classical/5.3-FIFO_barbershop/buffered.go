package main

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

const n = 3

var waitingRoom = make(chan int, n)
var barberRoom = make(chan int)
var wg sync.WaitGroup

func customer(c int) {
	time.Sleep(1e6)
	select {
	case barberRoom <- c:
		log.Print("customer ", c, " happy to find barber free")
	default:
		select {
		case waitingRoom <- c:
			log.Print("customer ", c, " waits")
		default:
			log.Print("customer ", c, " finds shop full, leaves")
			wg.Done()
		}
	}
}

func barber() {
	var c int
	for {
		log.Print("barber sleeping")
		select {
		case c = <-barberRoom:
		case c = <-waitingRoom:
		}
		log.Print("barber wakes and takes customer ", c)
		time.Sleep(1e8)
		wg.Done()
	awake:
		for {
			select {
			case c = <-waitingRoom:
				log.Print("barber takes waiting customer ", c)
				time.Sleep(1e8)
				wg.Done()
			default:
				break awake
			}
		}
	}
}

const nCust = 6

func main() {
	wg.Add(nCust)
	go barber()
	for c := 1; c <= nCust; c++ {
		go customer(c)
		time.Sleep(time.Duration(rand.Intn(1e7)))
	}
	wg.Wait()
}
