package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/konimarti/pipeline"
	"github.com/konimarti/pipeline/filters"
	"github.com/konimarti/pipeline/observer"
)

func main() {
	// define function
	var counter int
	norm := func() interface{} {
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
	monitor := pipeline.NewFromFunc(
		filters.NewChain(
			&filters.Stddev{Window: 20},
			&filters.Print{Writer: os.Stdout, Prefix: "Std Dev:"},
			&filters.AboveFloat64{1.4},
			&filters.Mute{Period: 2 * time.Second},
		),
		norm,
		500*time.Millisecond,
	)
	defer monitor.Close()

	// subscribers
	subscriber(1, monitor)
}

func subscriber(id int, monitor observer.Observer) {
	sub := monitor.Subscribe()
	for i := 0; i < 40; i++ {
		<-sub.Event()
		fmt.Printf("Std Dev: %2.4f -- Anomaly detected!\n", sub.Value().(float64))
		sub.Next()
	}
}
