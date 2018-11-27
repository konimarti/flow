 # Observer in Go

[![License](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://github.com/konimarti/observer/blob/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/konimarti/observer?status.svg)](https://godoc.org/github.com/konimarti/observer)
[![goreportcard](https://goreportcard.com/badge/github.com/konimarti/observer)](https://goreportcard.com/report/github.com/konimarti/observer)

Stream-processing Observer pattern for Golang with a modular notification behavior based on filters.

```go get github.com/konimarti/observer```

## Implementation Notes
Two types of observers are implemented which are suitable for different use cases:
* Channel-based observers accept new values through a ```chan interface{}``` channel, and
* Function-based observers collect new values in regular intervals from a ```func() interface{}``` function.

Channel-based observers are suitable in cases where we have control over the code and receive specific events. 
Function-based observer can monitor any object or state of resources (i.e. OPC servers without call-backs).

The filters control the behavior of the observer, i.e. they determine when and what value should be sent to the subscribers.  
This allows for a large flexibility and covers specific use cases by writing user-defined filters.
The following filters are currently implemented in this package:
- ```None{}```: No filter is applied. All values are send to the observers unfilitered and unprocessed.
- ```OnChange{}```: Notifies when the value changes.
- ```OnValue{value}```: Notifies when the new value matches the initial value.
- ```AboveFloat64{threshold}```: Notifies when a new float64 is above the initial float64 threshold.
- ```BelowFloat64{threshold}```: Notifies when a new float64 is below the initial float64 threshold.
- ```MovingAverage{Size}```: Calculates the moving average over a certain sample size and send the current moving mean to all subscribers.

## Stream processing: Anomaly detection and user-defined filters
An example for anomaly detection in streams using user-defined filters can be found [here](http://github.com/konimarti/observer/tree/master/example/anomaly_detection.go).

## Usage

* To get a channel-based observer:
```go
// define channel
ch := make(chan interface{})

// define filter
filter := filters.OnChange{}

// create observer
obs := observer.NewFromChan(&filter, chan interface{})

// publish new data to channel ch
// ch <- ..
```

* To get a function-based observer:
```go
// define a function that returns the values
fn := func() interface{} {
	return rand.Float64()
}

// define filter
filter := filters.OnChange{}

// create observer
obs := observer.NewFromFunc(&filter, fn, 1 * time.Second)
```

* Subscribers can subscribe to an observer and receive events that are triggered by the filter:
```go
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
```go
type Observer interface {
	Notify(interface{})
	Subscribe() Subscriber
	Close()
}
```

* Subscriber interface:
```go
type Subscriber interface {
	Value() interface{}
	Event() chan struct{}
	Next()
}
```

* Filter interface:
```go
type Filter interface {
	Check(interface{}) bool
	Update(interface{}) interface{}
}
```

## More Examples

See different examples [here](http://github.com/konimarti/observer/tree/master/example).

## Credits

The design of this observer implementation was inspired by [go-observer](http://github.com/imkira/go-observer).

## Disclaimer

This package is still work-in-progress. Interfaces might still change substantially. It is also not recommended to use it in a production environment at this point.





