package observer

//control struct is used to shut down an observer gracefully.
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

//NewControl creates a new control structure for graceful closing
//of the observer run loop
func NewControl() control {
	return control{C: make(chan bool), D: make(chan bool)}
}
