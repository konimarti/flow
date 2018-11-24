package observer

import (
	"io"
	"time"
)

type ValueFunc func() interface{}

//Observer containts the trigger and data gatherer
type observer struct {
	trigger   Trigger
	fn        ValueFunc
	observers []chan interface{}
	closing   []*control
}

//NewObserver creates a new observer struct
func NewObserver(t Trigger, f ValueFunc) *observer {
	return &observer{trigger: t, fn: f, observers: make([]chan interface{}, 0), closing: make([]*control, 0)}
}

//Notify sends out the current value in the observer channel
func (o *observer) Notify(value interface{}) {
	for _, observer := range o.observers {
		select {
		case <-observer:
		default:
		}
		observer <- value
	}
}

func (o *observer) Channel() chan interface{} {
	observer := make(chan interface{}, 1)
	o.observers = append(o.observers, observer)
	return observer
}

func (o *observer) Unsubscribe(ch chan interface{}) {
	// not implemented yet
}

func (o *observer) Observe(refresh time.Duration) io.Closer {

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

	return control
}

//Close closes all the observers channels
func (o *observer) Close() {
	for _, control := range o.closing {
		control.Close()
	}
}
