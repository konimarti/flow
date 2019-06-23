 # Stream processing flow in Go

[![License](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://github.com/konimarti/flow/blob/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/konimarti/flow?status.svg)](https://godoc.org/github.com/konimarti/flow)
[![goreportcard](https://goreportcard.com/badge/github.com/konimarti/flow)](https://goreportcard.com/report/github.com/konimarti/flow)

Stream processing in Golang with a modular notification behavior based on filters.

```go get github.com/konimarti/flow```

## Usage

1. Define a data stream source: 
```go
yourSource := flow.Func{
	func() interface{} {
		return rand.Float64()
	},
	500 * time.Millisecond,
}
```

2. Define the stream data processing (i.e. your filters):
```go
yourFilters := filters.AboveFloat64{0.5}
```

3. Create your flow and subscribe to the results:
```go
yourFlow := flow.New(
	yourFilters,
	yourSource,
)
results := yourFlow.Subscribe()
for {
	<-results.C()
	fmt.Println(results.Value())
}
```

## Example

* Apply a low-pass filter (exponential smoothing) to a sequence of random numbers between 0 and 1:

```go
flow := flow.New(
	&filters.Lowpass{A: 0.1}, 
	&flow.Func{ 
		func(){ return rand.Float64() },
		500 * time.Millisecond,
	},
)
```

* Generate a stream of random numbers from a standard normal, calculate moving average, print average and check if average is above or below 0.5 and -0.5, respectively. If so, then notify the subscribers.

```go
// define stream processor (flow) that returns an observer
yourFlow := flow.New(
	filters.NewChain(
		&filters.MovingAverage{Window: 10},
		&filters.Print{Writer: os.Stdout, Prefix: "Moving average:"},
		filters.NewSwitch(
			&filters.AboveFloat64{0.5},
			&filters.BelowFloat64{-0.5},
		),
	),
	&flow.Func{
		func() interface{} { return rand.NormFloat64() },
		500*time.Millisecond,
	},
)

// subscribe to flow and listen to events 
results := yourFlow.Subscribe()
for {
	<-results.C()
	fmt.Println("Notified:", results.Value())
}
```

## Description

Two types of flows are available that are suitable for different use cases:
* Channel-based flows accept new values through a ```chan interface{}``` channel, and
* Function-based flows collect new values in regular intervals from a ```func() interface{}``` function.

Channel-based flows are useful in cases where we receive specific events. 
Function-based flows can be used to monitor any object or state of resources 
(i.e. reading data from [OPC](http://github.com/konimarti/opc), HTTP requests, etc.).

* To get a channel-based flow:
```go
// define channel 
ch := make(chan interface{})

// define your filter(s) for data processing
filter := filters.OnChange{}

// create your flow
yourChannelFlow := flow.New(
	&filter, 
	&flow.Chan{ch},
)

// publish new data to channel ch
// ch <- ..
```

* To get a function-based flow:
```go
// define a function that returns the values
fn := func() interface{} {
	return rand.Float64()
}

// define your filter(s)
filter := filters.OnChange{}

// create your flow
yourFunctionFlow := flow.New(
	&filter, 
	&flow.Func{
		fn, 
		1 * time.Second,
	},
)
```

* You can subscribe to a flow in order to receive the results from the filter(s):
```go
results := yourFlow.Subscribe()
for {
	// wait for a result
	<-results.C()

	// get value from the filters
	results.Value()
}

```

## Filters

The filters control the behavior of the observer, i.e. they determine when and what values should be sent to the subscribers.  

### Available filters out-of-the-box

The following filters are currently implemented in this package:
* Notification filters:
  - ```None{}```: No filter is applied. All values are sent to the observers unfilitered and unprocessed.
  - ```Sink{}```: Blocks the flow of data. No values are sent to the observers.
  - ```Mute{Duration time.Duration}```: Mute shuts down all notifications after an event for a specific duration.
  - ```OnChange{}```: Notifies when the value changes.
  - ```OnValue{value float64}```: Notifies when the new value matches the defined value at initialization. 
  - ```AboveFloat64{threshold float64}```: Notifies when a new float64 is above the pre-defined float64 threshold.
  - ```BelowFloat64{threshold float64}```: Notifies when a new float64 is below the pre-defined float64 threshold.
  - ```Sigma{Window int, Factor float64}```: Sigma checks if the incoming value is a certain multiple (=factor) of standard deviations away from the mean.

* Stream-processing filters:
  - ```MovingAverage{Window int}```: Calculates the moving average over a certain sample size and sends the current mean to all subscribers.
  - ```StdDev{Window int}```: Calculates the standard deviation over a certain sample size and sends the current standard deviation to all subscribers.
  - ```LowPass{A float64}```: Performs low-pass filtering on the input data (exponential smoothing) with the smoothing factor A. 

### User-defined filters

User-defined filters can easily be created: Define your struct and embed the ```filters.Model```. You can then customize one or both of the interface functions. 
The ```filters.None``` is implemented by creating an empty struct and just embedding ```filters.Model```, for example.

The following user-defined filter expects a float64 value and multiplies it with a pre-defined factor:
```go
type Multiply struct {
	filters.Model
	Factor float64
}

func (m *Multiply) Update(v interface{}) interface{} {
	return v.(float64) * m.Factor
}
```

### Logical structures

Filters can be chained together using ```filters.NewChain(Filter1, Filter2, ...)```. 

To adjust the notification behavior, the ```filters.NewSwitch``` function can be useful, especially in cases when you want 
to monitor a value that needs to remain within a certain range ("deadband").

See [this example](http://github.com/konimarti/flow/tree/master/example/chain.go) for more information on logical structures 

### A stream-processing use case: Anomaly detection 

An anomaly detection example for streams with an user-defined filter based on Lytics' [Anomalyzer](http://github.com/lytics/anomalyzer) 
can be found [here](http://github.com/konimarti/flow/tree/master/example/anomaly.go).

## More examples

Check out the examples [here](http://github.com/konimarti/flow/tree/master/example).

## Credits

This software package has been developed for and is in production at [Kalkfabrik Netstal](http://www.kfn.ch/en).
The design of the observer implementation was inspired by [go-observer](http://github.com/imkira/go-observer).

## Disclaimer

This package is still work-in-progress. Interfaces might still change substantially. It is also not recommended to use it in a production environment at this point.





