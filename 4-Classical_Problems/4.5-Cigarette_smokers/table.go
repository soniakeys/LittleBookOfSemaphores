package main

import (
	"log"
	"math/rand"
	"sync"
)

type item int

const (
	paper item = iota
	tobacco
	matches
)

var itemString = []string{"paper", "tobacco", "matches"}

type table [3]int

var (
	t table
	m sync.Mutex
)

var wg sync.WaitGroup

func smoker(has item, report chan table) {
	var want1, want2 item
	switch has {
	case paper:
		want1, want2 = tobacco, matches
	case tobacco:
		want1, want2 = matches, paper
	case matches:
		want1, want2 = paper, tobacco
	}
	for r := range report {
		if r[want1] > 0 && r[want2] > 0 {
			took := false
			var tr table
			m.Lock()
			if t[want1] > 0 && t[want2] > 0 {
				t[want1]--
				t[want2]--
				took = true
				tr = t
			}
			m.Unlock()
			if took {
				log.Printf("smoker with %s smokes.  (p t m: %d %d %d)",
					itemString[has], tr[paper], tr[tobacco], tr[matches])
			}
		}
	}
	wg.Done()
}

const nRounds = 20

func main() {
	var rc [3]chan table
	for i := range rc {
		rc[i] = make(chan table)
		go smoker(item(i), rc[i])
	}
	for r := 0; r < nRounds; r++ {
		i := rand.Intn(3)
		m.Lock()
		t[i]++
		log.Printf("dealer puts %s (p t m: %d %d %d)",
			itemString[i], t[paper], t[tobacco], t[matches])
		tr := t
		m.Unlock()
		for _, ch := range rc {
			ch <- tr
		}
	}
	wg.Add(3)
	for _, ch := range rc {
		close(ch)
	}
	wg.Wait()
}
