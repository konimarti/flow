package pipeline

import (
	"testing"
	"time"
)

func TestControl(t *testing.T) {
	control := newControl()
	go func() {
		for {
			select {
			case <-control.C:
				control.D <- true
				return
			}
		}
	}()

	time.Sleep(100 * time.Millisecond)

	//Close
	ch := make(chan bool, 1)
	go func() {
		control.Close()
		ch <- true
	}()

	for {
		select {
		case <-time.After(2 * time.Second):
			t.Error("Closing timed out.")
		case <-ch:
			return
		}
	}
}
