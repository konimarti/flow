package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/konimarti/flow"
	"github.com/konimarti/flow/filters"
	"github.com/konimarti/flow/observer"
)

func main() {
	// set up function
	var counter int
	fn := func() interface{} {
		val := rand.Intn(4)
		if val == 3 {
			counter++
		}
		fmt.Printf("Publishing [count = %d]: %v \n", counter, val)
		return val
	}

	// create channel flow and use OnValue trigger
	monitor := flow.New(&filters.OnValue{3}, &flow.Func{fn, 10 * time.Millisecond})
	defer monitor.Close()

	// subscribers
	var wg sync.WaitGroup
	wg.Add(2)

	go subscriber(1, monitor, &wg)
	go subscriber(2, monitor, &wg)

	wg.Wait()
}

func subscriber(id int, monitor observer.Observer, wg *sync.WaitGroup) {
	sub := monitor.Subscribe()
	for i := 1; i < 10; i++ {
		<-sub.Event()
		fmt.Printf("Subscriber id(%d) got notified [%d]: %v\n", id, i, sub.Value())
		sub.Next()
	}
	wg.Done()
}
