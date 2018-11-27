package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/lytics/anomalyzer"

	"github.com/konimarti/observer"
	"github.com/konimarti/observer/filters"
)

type TransportLayer struct {
	Value float64
	Prob  float64
}

type AnomDetectFilter struct {
	filters.Model
	analyzer *anomalyzer.Anomalyzer
}

func (a *AnomDetectFilter) Update(v interface{}) interface{} {
	value := v.(float64)
	return TransportLayer{Value: value, Prob: a.analyzer.Push(value)}
}

func main() {
	// define anomalzyer
	conf := &anomalyzer.AnomalyzerConf{
		UpperBound: 3.0,
		LowerBound: -3.0,
		ActiveSize: 1,
		NSeasons:   12,
		Methods:    []string{"fence", "highrank", "lowrank"},
	}
	anom, _ := anomalyzer.NewAnomalyzer(conf, []float64{})

	// define function
	norm := func() interface{} {
		var anomaly float64
		if rand.Float64() < 0.1 {
			anomaly = float64(rand.Intn(10))
			if rand.Float64() < 0.5 {
				anomaly = -anomaly
			}
		}
		return rand.NormFloat64() + anomaly
	}

	// define function-based observer
	monitor := observer.NewFromFunc(&AnomDetectFilter{analyzer: &anom}, norm, 500*time.Millisecond)
	defer monitor.Close()

	// subscriber
	sub := monitor.Subscribe()
	for {
		<-sub.Event()
		tl := sub.Value().(TransportLayer)
		fmt.Printf("Value %+3.3f is anomalous with probabilty %3.3f", tl.Value, tl.Prob)
		if tl.Prob > 0.9 {
			fmt.Printf(" -- Anomaly detected!")
		}
		fmt.Printf("\n")
		sub.Next()
	}
}
