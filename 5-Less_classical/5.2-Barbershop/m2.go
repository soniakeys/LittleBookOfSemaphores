package main

import (
	"log"
	"os"
	"sync"
	"time"
)

const shopCap = 4 // number of customers allowed in shop at the same time
const nCust = 6   // total number of customers that will come to the shop

var (
	inShop = 0        // number of customers currently in shop
	mutex  sync.Mutex // protects inShop
)

// customer sends his customer number to enter the barber room
var barberRoom = make(chan int)

// barber sends dummy value when he's done cutting
var cutDone = make(chan int)

// counts customers as they leave shop
var wg sync.WaitGroup

func customer(c int) {
	mutex.Lock()
	log.Print("customer ", c, " arrives, sees ", inShop,
		" customers in shop")
	if inShop == shopCap {
		mutex.Unlock()
		log.Print("customer ", c, " finds shop full, leaves")
		wg.Done()
		return
	}
	inShop++
	log.Print("customer ", c, " waits")
	mutex.Unlock()
	barberRoom <- c
	time.Sleep(1e6)
	log.Print("customer ", c, " getting hair cut")
	<-cutDone
	mutex.Lock()
	inShop--
	mutex.Unlock()
	log.Print("customer ", c, " leaves with fresh hair cut")
	wg.Done()
}

func main() {
	wg.Add(nCust)
	go func() {
		wg.Wait()
		os.Exit(0)
	}()
	go func() {
		for c := 1; c <= nCust; c++ {
			go customer(c)
		}
	}()
	for {
		log.Print("barber sleeping")
		c := <-barberRoom
		log.Printf("barber wakes, cuts customer %d's hair", c)
		cutDone <- 1
	}
}
