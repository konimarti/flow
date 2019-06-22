package pipeline

import (
	"github.com/konimarti/pipeline/filters"
	"github.com/konimarti/pipeline/observer"
)

//NewFromChan creates a new observer struct
func NewFromChan(nf filters.Filter, channel chan interface{}) observer.Observer {
	obs := observer.NewObserver()
	runChan(obs, nf, channel)
	return obs
}

//runChan starts the observer with interval and fn
func runChan(o observer.Observer, nf filters.Filter, ch chan interface{}) {
	go func() {
		for {
			select {
			case v := <-ch:
				if nf.Check(v) {
					o.Notify(nf.Update(v))
				}
			case <-o.Control().C:
				o.Control().D <- true
				return
			}
		}
	}()
}
