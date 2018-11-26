# Observer in Go
Flexible Observer pattern for Golang with a trigger-based notification behavior.

```go get github.com/konimarti/observer```

## Notes on the implementations
Two type of observers are implemented suitable for different use cases:
* Channel-based observers accept new values through a ```chan interface{}``` channel, and
* Function-based observers that collect new values in regular intervals from a ```func() interface{}``` function.

Channel-based observers are suitable in cases where we have control over the code and receive specific events. Function-based observer can monitor any object or state of resources (i.e. OPC servers without call-backs).

The triggers control the behavior of the observer implementation and determines when to notify the observers. 
This gives a large flexibility and covers specific use cases with user-defined triggers.

## Example with a channel-based observer

```
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/konimarti/observer"
)

func main() {
	// set up channel
	ch := make(chan interface{})

	// create channel-based observer and set an OnValue trigger.
	// The observer will send notifications every time the value 3
	// is send through the channel.
	monitor := observer.NewFromChannel(&observer.OnValue{3}, ch)
	defer monitor.Close()

	// synchronization
	var wg sync.WaitGroup

	// publishers
	wg.Add(2)

	go publisher(1, ch, &wg)
	go publisher(2, ch, &wg)

	// subscribers
	wg.Add(2)

	go subscriber(1, monitor, &wg)
	go subscriber(2, monitor, &wg)

	wg.Wait()
}

func publisher(id int, ch chan interface{}, wg *sync.WaitGroup) {
	var counter int
	for {
		val := rand.Intn(4)
		if val == 3 {
			counter++
		}
		ch <- val
		sleep := rand.Intn(2)

		fmt.Printf("Publisher %d sends: value = %v, counts = %d, sleeps = %d sec \n", id, val, counter, sleep)
		time.Sleep(time.Duration(sleep) * time.Second)

		if counter >= 5 {
			break
		}
	}
	wg.Done()
}

func subscriber(id int, monitor observer.Observer, wg *sync.WaitGroup) {
	sub := monitor.Subscribe()
	for i := 1; i < 10; i++ {
		<-sub.Event()
		fmt.Printf("Subscriber %d got notified: value = %v, counts = %d\n", id, sub.Value(), i)
		sub.Next()
	}
	wg.Done()
}
```

## Example with a function-based observer

```
package main

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/konimarti/observer"
)

func main() {
	start := time.Now()
	freq := 0.05

	// define sinus function 
	sinfct := func() interface{} {
		sec := float64(time.Since(start).Seconds())
		sin := math.Sin(2.0 * math.Pi * freq * sec)
		fmt.Printf("x = %2.2f, sin(x) = %2.4f\n", sec, sin)
		return sin
	}

	// create function-based observer and set an AboveFloat64 trigger to send a notification
	// everytime the sinus function returns a value greater than 0.9.
	// The sinus function is evaluated every second.
	monitor := observer.NewFromFunction(&observer.AboveFloat64{0.9}, sinfct, 1*time.Second)
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
	for {
		<-sub.Event()
		fmt.Printf("Subscriber id(%d) got notified: %2.4f\n", id, sub.Value().(float64))
		sub.Next()
	}
	wg.Done()
}
```

## Credits

The design of this observer implementation was inspired by [go-observer](http://github.com/imkira/go-observer).

## Disclaimer

This package is still work-in-progress. Interfaces might still change substantially. It is also not recommended to use it in a production environment at this point.





