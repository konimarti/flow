package observer

import (
	"time"
)

// ValueFunc is the type of function to pass to the observer intsance to retrieve the next value
type ValueFunc func() interface{}

//observer contains the information to get a value
type observerInterval struct {
	observerImpl
}

//NewIntervalObserver creates a new observer struct
func NewIntervalObserver(tr Trigger, f ValueFunc, refresh time.Duration) Observer {
	obs := observerInterval{
		observerImpl{
			trigger: tr,
			state:   newState(),
			closing: make([]*control, 0),
		},
	}
	obs.run(time.Tick(refresh), f)
	return &obs
}

//run starts the observer with interval and fn
func (o *observerInterval) run(c <-chan time.Time, fn ValueFunc) {

	control := NewControl()
	o.observerImpl.closing = append(o.observerImpl.closing, control)

	go func() {
		for {
			select {
			case <-c:
				if v := fn(); o.observerImpl.trigger.Fire(v) {
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
