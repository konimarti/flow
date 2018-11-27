package observer

import (
	"time"

	"github.com/konimarti/observer/filters"
)

// ValueFunc is the type of function to pass to the observer intsance to retrieve the next value
type ValueFunc func() interface{}

//observer contains the information to get a value
type observerFunction struct {
	observerImpl
}

//NewFromFunction creates a new observer struct
func NewFromFunction(nf filters.Filter, f ValueFunc, refresh time.Duration) Observer {
	obs := observerFunction{
		observerImpl{
			control: newControl(),
			filter:  nf,
			state:   newState(),
		},
	}
	obs.run(time.Tick(refresh), f)
	return &obs
}

//run starts the observer with interval and fn
func (o *observerFunction) run(c <-chan time.Time, fn ValueFunc) {

	go func() {
		for {
			select {
			case <-c:
				if v := fn(); o.observerImpl.filter.Check(v) {
					o.Notify(o.filter.Update(v))
				}
			case <-o.control.C:
				o.control.D <- true
				return
			}
		}
	}()
}
