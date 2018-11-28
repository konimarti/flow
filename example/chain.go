package main

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/konimarti/observer"
	"github.com/konimarti/observer/filters"
)

func main() {
	// define function
	norm := func() interface{} {
		val := rand.NormFloat64()
		//fmt.Println("Publishing", val)
		return val
	}

	// Monitor Moving Average over 10 samples and notifies subscribers,
	// when average is below -0.5 or above 0.5.
	// Also, print out moving average with every update.
	monitor := observer.NewFromFunc(
		filters.NewChain(
			&filters.MovingAverage{Size: 10},
			&filters.Print{Writer: os.Stdout, Prefix: "Moving average:"},
			filters.NewOr(
				&filters.AboveFloat64{0.5},
				&filters.BelowFloat64{-0.5},
			),
		),
		norm,
		500*time.Millisecond,
	)
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
