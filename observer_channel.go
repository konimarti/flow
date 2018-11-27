package observer

import "github.com/konimarti/observer/filters"

type observerChannel struct {
	observerImpl
}

//NewFromChannel creates a new observer struct
func NewFromChan(nf filters.Filter, channel chan interface{}) Observer {
	obs := observerChannel{
		observerImpl: observerImpl{
			control: newControl(),
			filter:  nf,
			state:   newState(),
		},
	}
	obs.run(channel)
	return &obs
}

//run starts the observer with interval and fn
func (o *observerChannel) run(ch chan interface{}) {

	go func() {
		for {
			select {
			case v := <-ch:
				if o.observerImpl.filter.Check(v) {
					o.Notify(o.filter.Update(v))
				}
			case <-o.control.C:
				o.control.D <- true
				return
			}
		}
	}()
}
