# Observer in Go

[![License](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://github.com/konimarti/observer/blob/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/konimarti/observer?status.svg)](https://godoc.org/github.com/konimarti/observer)
[![goreportcard](https://goreportcard.com/badge/github.com/konimarti/observer)](https://goreportcard.com/report/github.com/konimarti/observer)

Flexible Observer pattern for Golang with a modular notification behavior based on filters.

```go get github.com/konimarti/observer```

## Implementation Notes
Two types of observers are implemented which are suitable for different use cases:
* Channel-based observers accept new values through a ```chan interface{}``` channel, and
* Function-based observers that collect new values in regular intervals from a ```func() interface{}``` function.

Channel-based observers are suitable in cases where we have control over the code and receive specific events. 
Function-based observer can monitor any object or state of resources (i.e. OPC servers without call-backs).

The filters control the behavior of the observers, i.e. they determine what value would trigger the notification of the subscribed observers.  
This gives a large flexibility and covers specific use cases with user-defined filters.
The following notifiers are currently implemented in this package:
- OnChange: Notifies when the value changes.
- OnValue: Notifies when the new value matches the initialized value.
- AboveFloat64: Notifies when a new float64 is above the initialized float64 threshold.
- BelowFloat64: Notifies when a new float64 is below the initialized float64 threshold.
- MovingAverage: Calculates the moving average over a certain sample sizes and informs all observers with the current value.

## Usage

* To get a channel-based observer:
```
// define channel
ch := make(chan interface{})

// define filter
filter := filters.OnChange{}

// create observer
obs := observer.NewFromChannel(&filter, chan interface{})

// publish new data to channel ch
// ch <- ..
```

* To get a function-based observer:
```
// define a function that returns the values
fn := func() interface{} {
	return rand.Float64()
}

// define filter
filter := filters.OnChange{}

// create observer
obs := observer.NewFromChannel(&filter, fn, 1 * time.Second)
```

* Subscribers can subscribe to a observer and receive events that are triggered by the filter used in the observer:
```
// subscribers
subscriber := obs.Subscribe()
for {
	// wait for event
	<-subscriber.Event()
	// get value that triggered event
	subscriber.Value()
	// advance to next
	subscriber.Next()
}

```

## Interfaces

* Observer interface:
```
type Observer interface {
	Notify(interface{})
	Subscribe() Subscriber
	Close()
}
```


* Filters interface:
```
type Filter interface {
	Check(interface{}) bool
	Update(interface{}) interface{}
}
```

## Examples

See different examples [here](https://github.com/konimarti/observer/tree/master/example).

## Credits

The design of this observer implementation was inspired by [go-observer](http://github.com/imkira/go-observer).

## Disclaimer

This package is still work-in-progress. Interfaces might still change substantially. It is also not recommended to use it in a production environment at this point.





