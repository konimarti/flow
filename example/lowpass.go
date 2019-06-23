package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/konimarti/flow"
	"github.com/konimarti/flow/filters"
)

func main() {
	// random numbers
	fn := func() interface{} {
		v := rand.Float64()
		fmt.Println("raw:", v)
		return v
	}

	// apply a low pass filter (exponential smoothing) to a sequency of random numbers between 0 and 1
	observer := flow.New(
		&filters.Lowpass{A: 0.1}, 
		&flow.Func{fn, 500 * time.Millisecond}
	)
	defer observer.Close()

	// subscribers
	sub := observer.Subscribe()
	for {
		<-sub.Event()
		fmt.Printf("filtered: %v\n", sub.Value())
		sub.Next()
	}
}
