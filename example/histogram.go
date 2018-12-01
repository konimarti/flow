package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/codahale/hdrhistogram"
	"github.com/konimarti/observer"
	"github.com/konimarti/observer/filters"
)

//Hist is a user-defined filter.
//Wrapper for hdrhistogram.
type HistFilter struct {
	hist *hdrhistogram.Histogram
	filters.Model
}

//Update expects int64 values and adds it to histogram.
//Returns the 99 percentile as int64.
func (h *HistFilter) Update(v interface{}) interface{} {
	h.hist.RecordValue(v.(int64))
	return h.hist.ValueAtQuantile(99.0)
}

func main() {
	// Exponential generator for int64
	fn := func() interface{} {
		return int64(math.Round(rand.ExpFloat64() * 100.0))
	}

	hist := hdrhistogram.New(0, 1000, 5)

	// define function-based observer
	monitor := observer.NewFromFunc(
		filters.NewChain(
			&HistFilter{hist: hist},
			&filters.Mute{Period: 1 * time.Second},
		),
		fn,
		1*time.Millisecond,
	)

	defer monitor.Close()

	// subscriber
	sub := monitor.Subscribe()
	for {
		<-sub.Event()
		fmt.Printf("Percentile: %v", sub.Value())
		fmt.Printf("\n")
		sub.Next()
	}
}
