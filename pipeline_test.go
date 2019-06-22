package pipeline_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/konimarti/pipeline"
	"github.com/konimarti/pipeline/filters"
)

var config = []struct {
	Values []interface{}
	Want   interface{}
}{
	{
		Values: []interface{}{1.1, 1.1, 1.1, 2.1, 1.1},
		Want:   2.1,
	},
	{
		Values: []interface{}{1, 1, 1, 2, 1},
		Want:   2,
	},
	{
		Values: []interface{}{"hello", "hello", "hello", "world", "hello"},
		Want:   "world",
	},
}

var observers = []struct {
	Name   string
	TrFunc func(v interface{}) filters.Filter
}{
	{
		Name: "OnChange",
		TrFunc: func(v interface{}) filters.Filter {
			return &filters.OnChange{v}
		},
	},
	{
		Name: "OnValue",
		TrFunc: func(v interface{}) filters.Filter {
			return &filters.OnValue{v}
		},
	},
}

func TestChannelObservers(t *testing.T) {

	for _, cfg := range config {
		for _, observerCfg := range observers {

			// prepare test
			startC := make(chan bool, 1)
			ch := make(chan interface{}, 1)
			values := cfg.Values
			go func() {
				<-startC
				for _, v := range values {
					ch <- v
					time.Sleep(10 * time.Millisecond)
				}
			}()

			// get start value
			var start interface{}
			switch observerCfg.Name {
			case "OnChange":
				start = values[0]
			case "OnValue":
				start = cfg.Want
			}

			// create observer
			observer := pipeline.New(observerCfg.TrFunc(start), &pipeline.Chan{ch})
			subscriber := observer.Subscribe()
			startC <- true

			// run test
			select {
			case <-time.After(1 * time.Second):
				str := fmt.Sprintf("%s: Timed out waiting for channel.", observerCfg.Name)
				t.Fatal(str)
			case <-subscriber.Event():
				received := subscriber.Value()
				if received != cfg.Want {
					str := fmt.Sprintf("%s: Got %v. Expected %v", observerCfg.Name, received, cfg.Want)
					t.Fatal(str)
				}
				subscriber.Next()
			}

			// close
			observer.Close()
		}
	}
}

func TestIntervalObservers(t *testing.T) {

	refresh := 10 * time.Millisecond

	for _, cfg := range config {
		for _, observerCfg := range observers {

			// prepare test
			values := cfg.Values
			var index int
			fn := func() interface{} {
				if index > len(values) {
					t.Error("Ran out of values.")
				}
				v := values[index]
				index++
				return v
			}

			// get start value
			var start interface{}
			switch observerCfg.Name {
			case "OnChange":
				start = values[0]
			case "OnValue":
				start = cfg.Want
			}

			// create observer
			observer := pipeline.New(observerCfg.TrFunc(start), &pipeline.Func{fn, refresh})
			subscriber := observer.Subscribe()
			// run test
			select {
			case <-time.After(1 * time.Second):
				str := fmt.Sprintf("%s: Timed out waiting for channel.", observerCfg.Name)
				t.Fatal(str)
			case <-subscriber.Event():
				received := subscriber.Value()
				if received != cfg.Want {
					str := fmt.Sprintf("%s: Got %v. Expected %v", observerCfg.Name, received, cfg.Want)
					t.Fatal(str)
				}
				subscriber.Next()
			}

			// close
			observer.Close()
		}
	}
}
