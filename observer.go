package pipeline

import (
	"sync"
)

type Observer interface {
	Notify(interface{})
	Subscribe() Subscriber
	Control() control
	Close()
}

//NewObserver
func NewObserver() Observer {
	o := observerI{
		control: newControl(),
		state:   newState(),
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
	next := newState()
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

// state
type state struct {
	C     chan struct{}
	Value interface{}
	Next  *state
}

func newState() *state {
	return &state{C: make(chan struct{})}
}
