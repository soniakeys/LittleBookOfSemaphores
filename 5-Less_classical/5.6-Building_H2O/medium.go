package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	nO = 10
	nH = 2 * nO
)

var medium = make([]interface{}, nO+nH)
var wg sync.WaitGroup

func oxygen(x int, bond chan int) {
	defer wg.Done()
	if x == 0 || x == len(medium)-1 {
		return
	}
	var neighbor1, neighbor2 chan int
	var ok bool
	if neighbor1, ok = medium[x-1].(chan int); !ok {
		return
	}
	if neighbor2, ok = medium[x+1].(chan int); !ok {
		return
	}
	select {
	case <-neighbor1:
	default:
		return
	}
	select {
	case <-neighbor2:
	default:
		neighbor1 <- '1'
		return
	}
	fmt.Println("bonding at", x)
	bond <- x
}

func main() {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < nH; i++ {
		ch := make(chan int, 1)
		ch <- 'H'
		medium[i] = ch
	}
	for len(medium) > 0 {
		rand.Shuffle(len(medium), func(i, j int) {
			medium[i], medium[j] = medium[j], medium[i]
		})
		for _, a := range medium {
			if a == nil {
				fmt.Print("O ")
			} else {
				fmt.Print("H ")
			}
		}
		fmt.Println()
		bond := make(chan int, nO)
		for i, atom := range medium {
			if atom == nil {
				wg.Add(1)
				go oxygen(i, bond)
			}
		}
		wg.Wait()
		close(bond)
		var bonded []int
		for x := range bond {
			// fixup x by already bonded atoms
			for _, b := range bonded {
				if x > b {
					x -= 3
				}
			}
			bonded = append(bonded, x)
			copy(medium[x-1:], medium[x+2:])
			medium = medium[:len(medium)-3]
		}
	}
}
