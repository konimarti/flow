package observer

import "sync"

// observer interface
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
	sync.RWMutex
	trigger Trigger
	state   *state // tip
	closing []*control
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

func (o *observerImpl) Subscribe() Subscriber {
	o.RLock()
	defer o.RUnlock()
	return &subscriber{state: o.state}
}

//Close closes all the observers channels
func (o *observerImpl) Close() {
	// TODO Just one?
	for _, control := range o.closing {
		control.Close()
	}
}
