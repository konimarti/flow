package observer

import "sync"

// Observer defines the observe interface
type Observer interface {
	Notify(interface{})
	Subscribe() Subscriber
	Close()
}

// state
type state struct {
	C     chan interface{}
	Value interface{}
	Next  *state
}

func newState() *state {
	return &state{C: make(chan interface{})}
}

//observerImpl implements the observer interface.
//This should be used in type embedding.
type observerImpl struct {
	sync.RWMutex //embedded
	control      //embedded
	trigger      Trigger
	state        *state
}

//Notify sends out the current value in the observer channel
func (o *observerImpl) Notify(value interface{}) {
	o.Lock()
	defer o.Unlock()
	o.state.Value = value
	next := newState()
	o.state.Next = next
	close(o.state.C)
	o.state = o.state.Next
}

//Subscribe returns a new subscriber to access values and listens for events
func (o *observerImpl) Subscribe() Subscriber {
	o.RLock()
	defer o.RUnlock()
	return &subscriber{state: o.state}
}
