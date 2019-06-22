package pipeline

import (
	"time"

	"github.com/konimarti/pipeline/filters"
)

// ValueFunc is the type of function to pass to the observer intsance to retrieve the next value
type ValueFunc func() interface{}

//NewFromFunc creates a new observer struct
func NewFromFunc(nf filters.Filter, f ValueFunc, refresh time.Duration) Observer {
	obs := NewObserver()
	runFunc(obs, nf, time.Tick(refresh), f)
	return obs
}

//runFunc
func runFunc(o Observer, nf filters.Filter, c <-chan time.Time, fn ValueFunc) {

	go func() {
		for {
			select {
			case <-c:
				if v := fn(); nf.Check(v) {
					o.Notify(nf.Update(v))
				}
			case <-o.Control().C:
				o.Control().D <- true
				return
			}
		}
	}()
}