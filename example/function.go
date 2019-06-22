package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/konimarti/pipeline"
	"github.com/konimarti/pipeline/filters"
	"github.com/konimarti/pipeline/observer"
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

	// create channel pipeline and use OnValue trigger
	monitor := pipeline.NewFromFunc(&filters.OnValue{3}, fn, 10*time.Millisecond)
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
