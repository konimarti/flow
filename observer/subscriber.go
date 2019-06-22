package observer

// Subscriber describes the interface
// returned by subscribing to an observer
type Subscriber interface {
	Value() interface{}
	Event() chan struct{}
	Next()
}

type subscriber struct {
	state *state
}

func (s *subscriber) Value() interface{} {
	return s.state.Value
}

func (s *subscriber) Event() chan struct{} {
	return s.state.C
}

func (s *subscriber) Next() {
	s.state = s.state.Next
}
