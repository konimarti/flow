package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/konimarti/flow"
	"github.com/konimarti/flow/filters"
	"github.com/konimarti/flow/observer"
)

func main() {
	// define function
	var counter int
	fn := func() interface{} {
		val := rand.NormFloat64()
		counter++
		factor := 1.0
		if counter > 40 && counter < 60 {
			factor = 2.0
		}
		return val * factor
	}

	// Monitors the running standard deviation of a data stream
	// and notifies the subscribers when the value reaches a
	// threshold of 1.4.
	flow := flow.New(
		filters.NewChain(
			&filters.Stddev{Window: 20},
			&filters.Print{Writer: os.Stdout, Prefix: "Std Dev:"},
			&filters.AboveFloat64{1.4},
			&filters.Mute{Period: 2 * time.Second},
		),
		&flow.Func{
			fn,
			500 * time.Millisecond,
		},
	)
	defer flow.Close()

	// subscribers
	subscriber(1, flow)
}

func subscriber(id int, flow observer.Observer) {
	sub := flow.Subscribe()
	for i := 0; i < 40; i++ {
		<-sub.C()
		fmt.Printf("Std Dev: %2.4f -- Anomaly detected!\n", sub.Value().(float64))
	}
}
