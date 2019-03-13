package pipeline

//The control struct to shut down all the observers gracefully.
//It implements the io.Closer interface.
type control struct {
	C chan bool
	D chan bool
}

//Close function closes the channel and waits for the done channel.
func (c *control) Close() {
	if c.C != nil && c.D != nil {
		c.C <- true
		<-c.D
	}
}

//newControl creates a new control structure for graceful closing
//of the observer run loop
func newControl() control {
	return control{C: make(chan bool), D: make(chan bool)}
}
