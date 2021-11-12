package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type CookInfo struct {
	foodCooked     string
	waitForPartner chan bool
}

func cookFood(player string) <-chan *CookInfo {
	cookChannel := make(chan *CookInfo)
	wait := make(chan bool)
	go func() {
		defer close(cookChannel)
		defer close(wait)
		for i := 0; ; i++ {
			cookChannel <- &CookInfo{fmt.Sprintf("%s %s", player, "Done"), wait}
			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)

			<-wait
		}
	}()
	return cookChannel
}

func fanIn(channels []<-chan *CookInfo) <-chan *CookInfo {
	var wg sync.WaitGroup
	channelOut := make(chan *CookInfo)

	// Wait to close channel
	go func() {
		wg.Wait()
		close(channelOut)
	}()

	for _, channel := range channels {
		wg.Add(1)
		go func(localCh <-chan *CookInfo) {
			defer wg.Done()
			for msg := range localCh {
				channelOut <- msg
			}
		}(channel)
	}

	return channelOut
}

func main() {
	gameChannel := fanIn([]<-chan *CookInfo{cookFood("Player 1 : "), cookFood("Player 2 :")})
	for i := 0; i < 3; i++ {
		food1 := <-gameChannel
		fmt.Println(food1.foodCooked)

		food2 := <-gameChannel
		fmt.Println(food2.foodCooked)

		food1.waitForPartner <- true
		food2.waitForPartner <- true
	}
}
