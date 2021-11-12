package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

// Generator Pattern
func channelGenerator(items int) <-chan string {
	channel := make(chan string)
	go func() {
		defer close(channel)
		for i := 0; i < items; i++ {
			channel <- strconv.Itoa(i)
		}
	}()
	return channel
}

type channelActor func(<-chan string) <-chan string

func doubleWithRandom(inputCh <-chan string) <-chan string {
	channel := make(chan string)
	go func() {
		defer close(channel)
		for range inputCh {
			rand.Seed(time.Now().UnixNano())
			num := rand.Intn(100)
			channel <- fmt.Sprintf("%d * 2 = %d", num, num*2)
		}
	}()
	return channel
}

// Fan-In Pattern
/**
ch1 -->|
	   ----> channelOut
ch2 -->|
*/
func fanIn(channels []<-chan string) <-chan string {
	var wg sync.WaitGroup
	channelOut := make(chan string)

	// Wait to close channel
	go func() {
		wg.Wait()
		close(channelOut)
	}()

	for _, channel := range channels {
		wg.Add(1)
		go func(localCh <-chan string) {
			defer wg.Done()
			for msg := range localCh {
				channelOut <- msg
			}
		}(channel)
	}

	return channelOut
}

// Fan-Out Pattern
/**
			   |--> ch1
channelIn ---> |
			   |--> ch2
*/
func fanOut(channelIn <-chan string, channelNums int, fn channelActor) []<-chan string {
	channelSlice := make([]<-chan string, channelNums)
	for i := 0; i < channelNums; i++ {
		channelSlice[i] = fn(channelIn)
	}
	return channelSlice
}

func timeExecution(what string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", what, time.Since(start))
	}
}

func main() {
	timer := timeExecution("Channel Generator")
	inputChannel := channelGenerator(1000000)
	timer()

	timer = timeExecution("Fan out")
	middleChannels := fanOut(inputChannel, 100, doubleWithRandom)
	timer()

	timer = timeExecution("Fan in")
	outputChannel := fanIn(middleChannels)
	timer()

	timer = timeExecution("Get results")
	for range outputChannel {
	}
	timer()
}
