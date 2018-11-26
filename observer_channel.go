package observer

import "github.com/konimarti/observer/notifiers"

type observerChannel struct {
	observerImpl
}

//NewFromChannel creates a new observer struct
func NewFromChannel(nf notifiers.Notifier, channel chan interface{}) Observer {
	obs := observerChannel{
		observerImpl: observerImpl{
			control:  newControl(),
			notifier: nf,
			state:    newState(),
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
				if o.observerImpl.notifier.Check(v) {
					o.Notify(v)
					o.notifier.Update(v)
				}
			case <-o.control.C:
				o.control.D <- true
				return
			}
		}
	}()
}
