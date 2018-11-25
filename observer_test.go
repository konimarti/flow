package observer_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/konimarti/observer"
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
	TrFunc func(v interface{}) observer.Trigger
}{
	{
		Name: "OnChange",
		TrFunc: func(v interface{}) observer.Trigger {
			return &observer.OnChange{v}
		},
	},
	{
		Name: "OnValue",
		TrFunc: func(v interface{}) observer.Trigger {
			return &observer.OnValue{v}
		},
	},
}

func TestObservers(t *testing.T) {

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
			observer := observer.NewObserver(observerCfg.TrFunc(start), fn, refresh)

			// run test
			select {
			case <-time.After(1 * time.Second):
				str := fmt.Sprintf("%s: Timed out waiting for channel.", observerCfg.Name)
				t.Fatal(str)
			case received := <-observer.Subscribe():
				if received != cfg.Want {
					str := fmt.Sprintf("%s: Got %v. Expected %v", observerCfg.Name, received, cfg.Want)
					t.Fatal(str)
				}
			}

			// close
			observer.Close()
		}
	}
}
