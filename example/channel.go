package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/konimarti/observer"
)

func main() {
	// set up channel
	ch := make(chan interface{}, 0)

	// create channel observer and use OnValue trigger
	monitor := observer.NewChannelObserver(&observer.OnValue{3}, ch)
	defer monitor.Close()

	// publisher: random numbers to be added in irregular intervals
	go publisher(1, ch)
	go publisher(2, ch)

	// subscribers
	var wg sync.WaitGroup
	wg.Add(2)

	go subscriber(1, monitor, &wg)
	go subscriber(2, monitor, &wg)

	wg.Wait()
}

func publisher(id int, ch chan interface{}) {
	var counter int
	for {
		val := rand.Intn(4)
		if val == 3 {
			counter++
		}
		ch <- val
		sleep := rand.Intn(2)
		fmt.Printf("Publishing id(%d) [%d]: %v (sleep for %ds)\n", id, counter, val, sleep)
		time.Sleep(time.Duration(sleep) * time.Second)
	}
}

func subscriber(id int, monitor observer.Observer, wg *sync.WaitGroup) {
	sub := monitor.Subscribe()
	for i := 1; i < 10; i++ {
		<-sub.Event()
		fmt.Printf("Subscriber id(%d) got notified [%d]: %v\n", id, i, sub.Value())
		sub.Next()
	}
	wg.Done()
}
