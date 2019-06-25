package flow

import (
	"time"

	"github.com/konimarti/flow/filters"
	"github.com/konimarti/flow/observer"
)

//New return an Observer that receive the results of the flow
func New(nf filters.Filter, s Source) observer.Observer {
	return s.Run(nf)
}

//Source is the interface for input for the flow
type Source interface {
	Run(f filters.Filter) observer.Observer
}

//Func implements the Source interface and regularly calls a function
type Func struct {
	Fn      func() interface{}
	Refresh time.Duration
}

//Run calls the given function in regular intervals
func (f *Func) Run(nf filters.Filter) observer.Observer {
	o := observer.NewObserver()
	c := time.Tick(f.Refresh)
	go func() {
		for {
			select {
			case <-c:
				if v := f.Fn(); nf.Check(v) {
					o.Notify(nf.Update(v))
				}
			case <-o.Control().C:
				o.Control().D <- true
				return
			}
		}
	}()
	return o
}

//Chan implements the Source interface and provides the input for the flow
type Chan struct {
	Ch chan interface{}
}

//Run passed the channel data to the filters
func (c *Chan) Run(nf filters.Filter) observer.Observer {
	o := observer.NewObserver()
	go func() {
		for {
			select {
			case v := <-c.Ch:
				if nf.Check(v) {
					o.Notify(nf.Update(v))
				}
			case <-o.Control().C:
				o.Control().D <- true
				return
			}
		}
	}()
	return o
}
