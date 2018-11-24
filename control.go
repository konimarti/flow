package observer

//The control struct to shut down all the observers gracefully.
//It implements the io.Closer interface.
type control struct {
	C chan bool
	D chan bool
}

//Close function closes the channel and waits for the done channel.
func (c *control) Close() error {
	if c.C != nil && c.D != nil {
		c.C <- true
		<-c.D
	}
	return nil
}

func NewControl() *control {
	return &control{C: make(chan bool), D: make(chan bool)}
}
