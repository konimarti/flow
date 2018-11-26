package main

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/konimarti/observer"
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

	// create function-based observer and set an AboveFloat64 trigger to send a notification
	// everytime the sinus function returns a value greater than 0.9.
	// The sinus function is evaluated every second.
	monitor := observer.NewFromFunction(&observer.AboveFloat64{0.9}, sinfct, 1*time.Second)
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
	for i := 0; i < 20; i++ {
		<-sub.Event()
		fmt.Printf("Subscriber id(%d) got notified: %2.4f\n", id, sub.Value().(float64))
		sub.Next()
	}
	wg.Done()
}
