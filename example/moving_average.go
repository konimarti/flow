package main

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/konimarti/pipeline"
	"github.com/konimarti/pipeline/filters"
	"github.com/konimarti/pipeline/observer"
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

	// create function-based pipeline and set an MovingAverage filter to
	// calulcate the moving average; expected moving average = 0.0
	// The function is evaluated every second.
	monitor := pipeline.New(&filters.MovingAverage{Window: 20},
		&pipeline.Func{sinfct, 1 * time.Second})
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
	for i := 0; i < 1000; i++ {
		<-sub.Event()
		fmt.Printf("Subscriber id(%d) got notified: %2.4f\n", id, sub.Value().(float64))
		sub.Next()
	}
	wg.Done()
}
