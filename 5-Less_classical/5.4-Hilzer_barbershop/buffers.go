package main

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

const (
	nBarber     = 3
	couchCap    = 4
	shopCap     = 20
	standingCap = shopCap - nBarber - couchCap
)

var barber1 = make(chan int)
var barber2 = make(chan int)
var barber3 = make(chan int)
var couch = make(chan int, couchCap)
var standing = make(chan int, standingCap-1)
var wg sync.WaitGroup

const nCust = 40

func main() {
	rand.Seed(time.Now().UnixNano())
	wg.Add(nCust)
	go func() {
		for {
			c := <-standing // c holds one standing customer
			couch <- c      // send blocks until there's a space on the couch
			log.Print("customer ", c, " takes the free seat on the couch")
		}
	}()
	go barber(1, barber1)
	go barber(2, barber2)
	go barber(3, barber3)
	log.Print("(business starts slow)")
	n := nCust / 6
	c := 1
	for ; c <= n; c++ {
		go customer(c)
		time.Sleep(time.Duration(rand.Intn(2e8)))
	}
	log.Print("(business picks up)")
	for ; c <= n*3; c++ {
		go customer(c)
		time.Sleep(time.Duration(rand.Intn(5e7)))
	}
	log.Print("(very busy now)")
	for ; c <= nCust; c++ {
		go customer(c)
		time.Sleep(time.Duration(rand.Intn(8e6)))
	}
	wg.Wait()
}

func barber(b int, ch chan int) {
	var c int
	cut := func() {
		// time for haircut varies
		time.Sleep(time.Duration(5e7 + rand.Intn(5e7)))
		log.Print("customer ", c, " pays and leaves")
		wg.Done()
	}
	for {
		log.Print("barber ", b, " sleeping")
		select {
		case c = <-ch:
		case c = <-couch:
		}
		log.Print("barber ", b, " wakes and takes customer ", c)
		cut()
	awake:
		for {
			select {
			case c = <-couch:
				log.Print("barber ", b, " takes waiting customer ", c,
					" from the couch")
				cut()
			default:
				break awake
			}
		}
	}
}

func customer(c int) {
	time.Sleep(1e6) // customer spends a brief moment scoping things out
	select {
	case barber1 <- c:
		log.Print("customer ", c, " arrives and is happy to find a barber free")
	case barber2 <- c:
		log.Print("customer ", c, " arrives and is happy to find a barber free")
	case barber3 <- c:
		log.Print("customer ", c, " arrives and is happy to find a barber free")
	default:
		select {
		case couch <- c:
			log.Print("customer ", c, " arrives and takes a seat on the couch")
		default:
			select {
			case standing <- c:
				log.Print("customer ", c, " arrives and waits standing")
			default:
				log.Print("customer ", c, " arrives, finds shop full, leaves")
				wg.Done()
			}
		}
	}
}
