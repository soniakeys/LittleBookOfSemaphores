package main

import (
	"log"
	"math/rand"
	"sync"
)

var round sync.WaitGroup

func smoker(have string, want1, want2, take chan int) {
	var oneItem bool
watch:
	for {
		select {
		case <-want1:
			if oneItem {
				break watch
			}
			oneItem = true
			want1 = nil
		case <-want2:
			if oneItem {
				break watch
			}
			oneItem = true
			want2 = nil
		case <-take:
			return
		}
	}
	close(take)
	log.Println("smoker with", have, "smokes")
	round.Done()
}

const nRounds = 10

func main() {
	round.Add(nRounds)
	for r := 0; r < nRounds; r++ {
		paper := make(chan int)
		tobacco := make(chan int)
		match := make(chan int)
		take := make(chan int)
		go smoker("paper", tobacco, match, take)
		go smoker("tobacco", match, paper, take)
		go smoker("match", paper, tobacco, take)
		switch rand.Intn(3) {
		case 0:
			log.Println("dealer puts paper and tobacco")
			close(paper)
			close(tobacco)
		case 1:
			log.Println("dealer puts tobacco and match")
			close(tobacco)
			close(match)
		case 2:
			log.Println("dealer puts match and paper")
			close(match)
			close(paper)
		}
		<-take
	}
	round.Wait()
}
