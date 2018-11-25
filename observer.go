package observer

import (
	"time"
)

type ValueFunc func() interface{}

//observer contains the information to get a value
type observer struct {
	trigger   Trigger
	fn        ValueFunc
	observers []chan interface{}
	closing   []*control
}

//NewObserver creates a new observer struct
func NewObserver(tr Trigger, f ValueFunc, rf time.Duration) *observer {
	obs := observer{trigger: tr, fn: f, observers: make([]chan interface{}, 0), closing: make([]*control, 0)}
	obs.run(rf)
	return &obs
}

//Notify sends out the current value in the observer channel
func (o *observer) Notify(value interface{}) {
	for _, observer := range o.observers {
		select {
		case v := <-observer:
			// in case channel has been closed,
			// remove it from the list
			if v == nil {
				o.Unsubscribe(observer)
				continue
			}
		default:
		}
		observer <- value
	}
}

func (o *observer) Subscribe() chan interface{} {
	observer := make(chan interface{}, 1)
	o.observers = append(o.observers, observer)
	return observer
}

func (o *observer) Unsubscribe(ch chan interface{}) {
	// not sure if needed
}

//run starts the observer
func (o *observer) run(refresh time.Duration) {

	control := NewControl()
	o.closing = append(o.closing, control)

	c := time.Tick(refresh)

	go func() {
		for {
			select {
			case <-c:
				if v := o.fn(); o.trigger.Fire(v) {
					o.Notify(v)
					o.trigger.Update(v)
				}
			case <-control.C:
				control.D <- true
				return
			}
		}
	}()
}

//Close closes all the observers channels
func (o *observer) Close() {
	for _, control := range o.closing {
		control.Close()
	}
}
