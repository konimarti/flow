package main

import (
	"fmt"

	"github.com/konimarti/flow"
	"github.com/konimarti/flow/filters"
)

func main() {
	// set up channel
	ch := make(chan interface{})

	// stream strings through the channel to the filters
	input := []string{"Alabama", "Alaska", "Arizona", "Arkensas", "California", "Colorado"}
	done := make(chan struct{})
	go func() {
		for _, word := range input {
			fmt.Println("entering:", word)
			ch <- word
		}
		done <- struct{}{}
	}()

	// create channel-based flow and set an OnValue trigger.
	flow := flow.New(&filters.OnValue{"California"}, &flow.Chan{ch})
	defer flow.Close()

	// get the results of the stream processing
	results := flow.Subscribe()
	for {
		select {
		case <-results.C():
			fmt.Println("Found:", results.Value())
		case <-done:
			return
		}
	}
}
