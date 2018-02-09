package main

import (
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

var waitingRoomChair1 = make(chan int, 1)
var waitingRoomChair2 = make(chan int, 1)
var waitingRoomChair3 = make(chan int, 1)
var barberChair = make(chan int)
var wg sync.WaitGroup

func customerFunc(c int) {
	time.Sleep(1e6)
	select {
	case barberChair <- c:
		log.Print("customer ", c, " happy to find barber free")
	default:
		select {
		case waitingRoomChair1 <- c:
			log.Print("customer ", c, " waits")
		case waitingRoomChair2 <- c:
			log.Print("customer ", c, " waits")
		case waitingRoomChair3 <- c:
			log.Print("customer ", c, " waits")
		default:
			log.Print("customer ", c, " finds shop full, leaves")
			wg.Done()
		}
	}
}

const nCust = 6

func main() {
	wg.Add(nCust)
	go func() {
		wg.Wait()
		os.Exit(0)
	}()
	go func() {
		for c := 1; c <= nCust; c++ {
			go customerFunc(c)
			time.Sleep(time.Duration(rand.Intn(1e7)))
		}
	}()
	var c int
	for {
		log.Print("barber sleeping")
		select {
		case c = <-barberChair:
		case c = <-waitingRoomChair1:
		case c = <-waitingRoomChair2:
		case c = <-waitingRoomChair3:
		}
		log.Print("barber wakes and takes customer ", c)
		time.Sleep(1e8)
		wg.Done()
	awake:
		for {
			select {
			case c = <-waitingRoomChair1:
			case c = <-waitingRoomChair2:
			case c = <-waitingRoomChair3:
			default:
				break awake
			}
			log.Print("barber takes waiting customer ", c)
			time.Sleep(1e8)
			wg.Done()
		}
	}
}
