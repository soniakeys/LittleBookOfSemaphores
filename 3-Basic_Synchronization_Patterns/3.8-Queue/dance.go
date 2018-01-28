package main

import (
	"log"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())
	leaders := make(chan int)
	followers := make(chan int)
	danceEnds := time.After(time.Second)
	randInt := int(3e8)
	leaderArrives := time.After(time.Duration(rand.Intn(randInt)))
	followerArrives := time.After(time.Duration(rand.Intn(randInt)))
	var l, f int
	for {
		select {
		case <-danceEnds:
			log.Println("dance ends")
			return
		case <-leaderArrives:
			l++
			log.Println("leader", l, "arrives")
			select {
			case f := <-followers:
				log.Println("leader", l, "follower", f, "dance")
			default:
				go func(l int) { leaders <- l }(l)
			}
			leaderArrives = time.After(time.Duration(rand.Intn(randInt)))
		case <-followerArrives:
			f++
			log.Println("follower", f, "arrives")
			select {
			case l := <-leaders:
				log.Println("leader", l, "follower", f, "dance")
			default:
				go func(f int) { followers <- f }(f)
			}
			followerArrives = time.After(time.Duration(rand.Intn(randInt)))
		}
	}
}
