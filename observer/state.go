package observer

// state
type state struct {
	C     chan struct{}
	Value interface{}
	Next  *state
}

//NewState creats a new state
func NewState() *state {
	return &state{C: make(chan struct{})}
}
