package pipeline

import (
	"time"

	"github.com/konimarti/pipeline/filters"
	"github.com/konimarti/pipeline/observer"
)

//New
func New(nf filters.Filter, s Source) observer.Observer {
	return s.Run(nf)
}

//Source
type Source interface {
	Run(f filters.Filter) observer.Observer
}

//Func
type Func struct {
	Fn      func() interface{}
	Refresh time.Duration
}

//Run
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

//Chan
type Chan struct {
	Ch chan interface{}
}

//Run
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
