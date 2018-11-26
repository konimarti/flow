package observer_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/konimarti/observer"
)

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
			observer := observer.NewFromChannel(observerCfg.TrFunc(start), ch)
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
