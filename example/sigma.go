package main

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/konimarti/flow"
	"github.com/konimarti/flow/filters"
	"github.com/konimarti/flow/observer"
)

func main() {
	// define function
	fn := func() interface{} {
		val := rand.NormFloat64()
		//fmt.Println("Publishing", val)
		return val
	}

	// Monitor Moving Average over 10 samples and notifies subscribers,
	// when average is below -0.5 or above 0.5.
	// Also, print out moving average with every update.
	flow := flow.New(
		filters.NewChain(
			&filters.Print{Writer: os.Stdout, Prefix: "input"},
			&filters.Sigma{Window: 10, Factor: 2},
			&filters.Print{Writer: os.Stdout, Prefix: "output"},
		),
		&flow.Func{
			fn,
			500 * time.Millisecond,
		},
	)
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
	for i := 0; i < 5; i++ {
		<-sub.C()
		fmt.Printf("Subscriber id(%d) got notified: %2.4f\n", id, sub.Value().(float64))
	}
	wg.Done()
}
