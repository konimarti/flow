package observer

import (
	"sync"
)

//Observer defines the interface for the observer design patter
type Observer interface {
	Notify(interface{})
	Subscribe() Subscriber
	Control() control
	Close()
}

//NewObserver returns an implementation of the observer interface
func NewObserver() Observer {
	o := observerI{
		control: NewControl(),
		state:   NewState(),
	}
	return &o
}

//observer
type observerI struct {
	sync.RWMutex //embedded
	control      //embedded
	state        *state
}

//Notify sends out the current value in the observer channel
func (o *observerI) Notify(value interface{}) {
	o.Lock()
	defer o.Unlock()
	o.state.Value = value
	next := NewState()
	o.state.Next = next
	close(o.state.C)
	o.state = o.state.Next
}

//Subscribe returns a new subscriber to access values and listens for events
func (o *observerI) Subscribe() Subscriber {
	o.RLock()
	defer o.RUnlock()
	return &subscriber{state: o.state}
}

//Control
func (o *observerI) Control() control {
	return o.control
}
