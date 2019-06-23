package observer

// Subscriber describes the interface
// returned by subscribing to an observer
type Subscriber interface {
	C() chan struct{}
	Value() interface{}
}

type subscriber struct {
	state *state
}

//Value return the current value
func (s *subscriber) Value() interface{} {
	v := s.v()
	s.next()
	return v
}

//C returns a channel and signals if Value() can be called
func (s *subscriber) C() chan struct{} {
	return s.state.C
}

func (s *subscriber) v() interface{} {
	return s.state.Value
}

func (s *subscriber) next() {
	s.state = s.state.Next
}
