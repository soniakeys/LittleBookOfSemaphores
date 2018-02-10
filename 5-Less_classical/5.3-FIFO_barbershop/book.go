package main

import (
	"log"
	"sync"

	"github.com/soniakeys/LittleBookOfSemaphores/sem"
)

var (
	n            = 4
	customers    = 0
	mutex        = sem.NewChanSem(1, 1)
	customer     = sem.NewChanSem(0, 1)
	customerDone = sem.NewChanSem(0, 1)
	barberDone   = sem.NewChanSem(0, 1)
	queue        []sem.ChanSem
)

var wg sync.WaitGroup

func customerFunc(c int) {
	s := sem.NewChanSem(0, 1)
	mutex.Wait()
	log.Print("customer ", c, " arrives, sees ", customers,
		" customers in shop")
	if customers == n {
		mutex.Signal()
		log.Print("customer ", c, " finds shop full, leaves")
		wg.Done()
		balk()
	}
	customers++
	log.Print("customer ", c, " waits")
	queue = append(queue, s)
	mutex.Signal()

	customer.Signal()
	s.Wait()

	// getHairCut ()
	log.Print("customer ", c, " gets hair cut")

	customerDone.Signal()
	barberDone.Wait()

	mutex.Wait()
	customers--
	mutex.Signal()
	log.Print("customer ", c, " leaves with fresh hair cut")
	wg.Done()
}

func balk() {
	select {}
}

func barberFunc() {
	for {
		log.Print("barber sleeping")

		customer.Wait()
		mutex.Wait()
		s := queue[0]
		queue = queue[1:]

		mutex.Signal()
		s.Signal()

		// cutHair ()
		log.Print("barber cutting hair")

		customerDone.Wait()
		barberDone.Signal()
	}
}

const nCust = 6

func main() {
	wg.Add(nCust)
	go barberFunc()
	for i := 1; i <= nCust; i++ {
		go customerFunc(i)
	}
	wg.Wait()
}
