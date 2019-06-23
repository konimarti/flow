package main

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/konimarti/flow"
	"github.com/konimarti/flow/filters"
	"github.com/konimarti/flow/observer"
)

func main() {
	start := time.Now()
	freq := 0.05

	// define sinus function
	sinfct := func() interface{} {
		sec := float64(time.Since(start).Seconds())
		sin := math.Sin(2.0 * math.Pi * freq * sec)
		fmt.Printf("x = %2.2f, sin(x) = %2.4f\n", sec, sin)
		return sin
	}

	// create function-based flow and set an MovingAverage filter to
	// calulcate the moving average; expected moving average = 0.0
	// The function is evaluated every second.
	flow := flow.New(&filters.MovingAverage{Window: 20},
		&flow.Func{sinfct, 1 * time.Second})
	defer flow.Close()

	// subscribers
	var wg sync.WaitGroup
	wg.Add(2)

	go subscriber(1, flow, &wg)
	go subscriber(2, flow, &wg)

	wg.Wait()
}

func subscriber(id int, flow observer.Observer, wg *sync.WaitGroup) {
	sub := flow.Subscribe()
	for i := 0; i < 1000; i++ {
		<-sub.C()
		fmt.Printf("Subscriber id(%d) got notified: %2.4f\n", id, sub.Value().(float64))
	}
	wg.Done()
}
