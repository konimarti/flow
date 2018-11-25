package main

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/konimarti/observer"
)

func main() {
	// set up function
	start := time.Now()
	freq := 0.05
	fn := func() interface{} {
		sec := float64(time.Since(start).Seconds())
		sin := math.Sin(2.0 * math.Pi * freq * sec)
		fmt.Printf("t = %2.2f    sin(t) = %2.4f\n", sec, sin)
		return sin
	}

	// create channel observer and use OnValue trigger
	monitor := observer.NewIntervalObserver(&observer.AboveFloat64{0.9}, fn, 1*time.Second)
	defer monitor.Close()

	// subscribers
	var wg sync.WaitGroup
	wg.Add(2)

	go subscriber(1, monitor, &wg)
	go subscriber(2, monitor, &wg)

	wg.Wait()
}

func subscriber(id int, monitor observer.Observer, wg *sync.WaitGroup) {
	sub := monitor.Subscribe()
	for {
		<-sub.Event()
		fmt.Printf("Subscriber id(%d) got notified: %2.4f\n", id, sub.Value().(float64))
		sub.Next()
	}
	wg.Done()
}
