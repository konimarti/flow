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
	// define generator function
	fn := func() interface{} {
		return rand.NormFloat64()
	}

	// Flow consits of calculatingaa moving average with 10 samples,
	// printing out the average,
	// and checking if it is outside a certain boundary.
	flow := flow.New(
		filters.NewChain(
			&filters.MovingAverage{Window: 10},
			&filters.Print{Writer: os.Stdout, Prefix: "Moving average:"},
			filters.NewSwitch(
				&filters.AboveFloat64{0.5},
				&filters.BelowFloat64{-0.5},
			),
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
	for i := 0; i < 20; i++ {
		<-sub.C()
		fmt.Printf("Subscriber id(%d) got notified: %2.4f\n", id, sub.Value().(float64))
	}
	wg.Done()
}
