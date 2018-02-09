package main

import (
	"log"
	"os"
	"sync"
)

var (
	n            = 4
	customers    = 0
	mutex        sync.Mutex
	customer     = make(chan int)
	barber       = make(chan int)
	customerDone = make(chan int)
	barberDone   = make(chan int)
)

var wg sync.WaitGroup

func customerFunc(c int) {
	mutex.Lock()
	log.Print("customer ", c, " arrives, sees ", customers,
		" customers in shop")
	if customers == n {
		mutex.Unlock()
		log.Print("customer ", c, " finds shop full, leaves")
		wg.Done()
		return
	}
	customers++
	log.Print("customer ", c, " waits")
	mutex.Unlock()
	customer <- c
	<-barber
	log.Print("customer ", c, " gets hair cut")
	customerDone <- 1
	<-barberDone
	mutex.Lock()
	customers--
	mutex.Unlock()
	log.Print("customer ", c, " leaves with fresh hair cut")
	wg.Done()
}

const nCust = 6

func main() {
	wg.Add(nCust)
	go func() {
		wg.Wait()
		os.Exit(0)
	}()
	for i := 1; i <= nCust; i++ {
		go customerFunc(i)
	}
	for {
		log.Print("barber sleeping")
		c := <-customer
		log.Print("barber wakes and seats customer ", c)
		barber <- 1
		log.Printf("barber cutting customer %d's hair", c)
		<-customerDone
		barberDone <- 1
	}
}
