package main

import (
	"log"
	"sync"

	"github.com/soniakeys/LittleBookOfSemaphores/sem"
)

var (
	n         = 20
	customers = 0
	mutex     = sem.NewChanSem(1, 1)
	mutex2    = sem.NewChanSem(1, 1)
	sofa      = sem.NewChanSem(4, 4)
	customer1 = sem.NewChanSem(0, 1)
	customer2 = sem.NewChanSem(0, 1)
	payment   = sem.NewChanSem(0, 1)
	receipt   = sem.NewChanSem(0, 1)
	queue1    []sem.ChanSem
	queue2    []sem.ChanSem
)

var wg sync.WaitGroup

func customerFunc(c int) {
	s1 := sem.NewChanSem(0, 1)
	s2 := sem.NewChanSem(0, 1)
	mutex.Wait()
	log.Print("customer ", c, " arrives, sees ", customers,
		" customers in shop")
	if customers == n {
		mutex.Signal()
		log.Print("customer ", c, " finds shop full, leaves")
		wg.Done()
		balk()
	}
	log.Print("customer ", c, " enters waiting area")
	customers++
	queue1 = append(queue1, s1)
	mutex.Signal()

	// enterShop ()
	customer1.Signal()
	s1.Wait()

	sofa.Wait()
	// sitOnSofa ()
	s1.Signal()
	mutex2.Wait()
	log.Print("customer ", c, " sits on sofa")
	queue2 = append(queue2, s2)
	mutex2.Signal()
	customer2.Signal()
	s2.Wait()
	sofa.Signal()

	// sitInBarberChair ()
	log.Print("customer ", c, " gets hair cut")

	// pay ()
	mutex.Wait()
	log.Print("customer ", c, " pays")
	payment.Signal()
	receipt.Wait()
	customers--
	mutex.Signal()
	log.Print("customer ", c, " leaves with fresh hair cut")
	wg.Done()
}

func balk() {
	select {}
}

func barberFunc(b int) {
	for {
		log.Print("barber ", b, " sleeping")
		customer1.Wait()
		mutex.Wait()
		s := queue1[0]
		queue1 = queue1[1:]
		s.Signal()
		s.Wait()
		mutex.Signal()

		customer2.Wait()
		mutex2.Wait()
		s = queue2[0]
		queue2 = queue2[1:]
		s.Signal()
		mutex2.Signal()

		// cutHair ()
		log.Print("barber ", b, " cutting hair")

		payment.Wait()
		// acceptPayment ()
		log.Print("barber ", b, " accepts payment")
		receipt.Signal()
	}
}

const nCust = 25

func main() {
	wg.Add(nCust)
	for b := 1; b <= 3; b++ {
		go barberFunc(b)
	}
	for c := 1; c <= nCust; c++ {
		go customerFunc(c)
	}
	wg.Wait()
}
