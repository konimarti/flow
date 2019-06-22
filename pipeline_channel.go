package pipeline

import "github.com/konimarti/pipeline/filters"

//NewFromChan creates a new observer struct
func NewFromChan(nf filters.Filter, channel chan interface{}) Observer {
	obs := NewObserver()
	runChan(obs, nf, channel)
	return obs
}

//runChan starts the observer with interval and fn
func runChan(o Observer, nf filters.Filter, ch chan interface{}) {
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
