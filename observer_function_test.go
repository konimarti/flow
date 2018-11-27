package observer_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/konimarti/observer"
)

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
			observer := observer.NewFromFunc(observerCfg.TrFunc(start), fn, refresh)
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
