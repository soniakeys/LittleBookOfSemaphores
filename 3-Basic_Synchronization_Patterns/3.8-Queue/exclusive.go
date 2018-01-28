package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())
	const pairs = 5
	const dancers = pairs * 2
	var ready [dancers]bool // array representation of a set
	for _, n := range rand.Perm(dancers) {
		if n < pairs {
			pair := n
			l := pair
			f := n + pairs
			fmt.Println("leader", pair, "ready")
			if ready[f] {
				fmt.Println("pair", pair, "dances")
			} else {
				ready[l] = true
			}
		} else {
			pair := n - pairs
			l := pair
			f := n
			fmt.Println("follower", pair, "ready")
			if ready[l] {
				fmt.Println("pair", pair, "dances")
			} else {
				ready[f] = true
			}
		}
	}
}
