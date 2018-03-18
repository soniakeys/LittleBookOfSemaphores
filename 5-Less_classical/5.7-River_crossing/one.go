package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	nHackers = 20
	nSerfs   = 20
)

func main() {
	rand.Seed(time.Now().UnixNano())
	arrivals := make([]byte, nHackers+nSerfs)
	for i := 0; i < nHackers; i++ {
		arrivals[i] = 'h'
	}
	rand.Shuffle(len(arrivals), func(i, j int) {
		arrivals[i], arrivals[j] = arrivals[j], arrivals[i]
	})
	var hWaiting, sWaiting int
	for _, a := range arrivals {
		if a == 'h' {
			hWaiting++
		} else {
			sWaiting++
		}
		switch {
		case hWaiting == 4:
			hWaiting -= 4
			fmt.Println("4 hackers cross")
		case sWaiting == 4:
			sWaiting -= 4
			fmt.Println("4 serfs cross")
		case hWaiting >= 2 && sWaiting >= 2:
			hWaiting -= 2
			sWaiting -= 2
			fmt.Println("2 hackers, 2 serfs cross")
		}
	}
}
