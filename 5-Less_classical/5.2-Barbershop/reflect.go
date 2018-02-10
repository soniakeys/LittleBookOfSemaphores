package main

import (
	"log"
	"math/rand"
	"reflect"
	"sync"
	"time"
)

const nWRChairs = 3

var wrChairs = make([]reflect.Value, nWRChairs)
var barberChair = make(chan int)
var wg sync.WaitGroup

const nCust = 6

func main() {
	wg.Add(nCust)
	for i := range wrChairs {
		wrChairs[i] = reflect.MakeChan(reflect.TypeOf(barberChair), 1)
	}
	go barber()
	for c := 1; c <= nCust; c++ {
		go customer(c)
		time.Sleep(time.Duration(rand.Intn(1e7)))
	}
	wg.Wait()
}

func barber() {
	sleepingCases := make([]reflect.SelectCase, nWRChairs+1)
	for i := 0; i < nWRChairs; i++ {
		sleepingCases[i] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: wrChairs[i],
		}
	}
	sleepingCases[nWRChairs] = reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(barberChair),
	}
	awakeCases := append([]reflect.SelectCase{}, sleepingCases...)
	awakeCases[nWRChairs] = reflect.SelectCase{
		Dir: reflect.SelectDefault,
	}
	for {
		log.Print("barber sleeping")
		_, recv, _ := reflect.Select(sleepingCases)
		log.Print("barber wakes and takes customer ", recv.Int())
		time.Sleep(1e8)
		wg.Done()
		// awake
		for {
			chosen, recv, _ := reflect.Select(awakeCases)
			if chosen == nWRChairs {
				break
			}
			log.Print("barber takes waiting customer ", recv.Int())
			time.Sleep(1e8)
			wg.Done()
		}
	}
}

func customer(c int) {
	time.Sleep(1e6)
	cVal := reflect.ValueOf(c)
	cases := make([]reflect.SelectCase, nWRChairs+1)
	for i := 0; i < nWRChairs; i++ {
		cases[i] = reflect.SelectCase{
			Dir:  reflect.SelectSend,
			Chan: wrChairs[i],
			Send: cVal,
		}
	}
	cases[nWRChairs].Dir = reflect.SelectDefault
	select {
	case barberChair <- c:
		log.Print("customer ", c, " happy to find barber free")
	default:
		if chosen, _, _ := reflect.Select(cases); chosen < nWRChairs {
			log.Print("customer ", c, " waits")
		} else {
			log.Print("customer ", c, " finds shop full, leaves")
			wg.Done()
		}
	}
}
