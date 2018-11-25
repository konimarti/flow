package observer

type observerChannel struct {
	observerImpl
}

//NewIntervalObserver creates a new observer struct
func NewChannelObserver(tr Trigger, channel chan interface{}) Observer {
	obs := observerChannel{
		observerImpl: observerImpl{
			trigger:   tr,
			observers: make([]chan interface{}, 0),
			closing:   make([]*control, 0),
		},
	}
	obs.run(ch)
	return &obs
}

//run starts the observer with interval and fn
func (o *observerChannel) run(ch chan interface{}) {

	control := NewControl()
	o.observerImpl.closing = append(o.observerImpl.closing, control)

	go func() {
		for {
			select {
			case v := <-ch:
				if o.observerImpl.trigger.Fire(v) {
					o.Notify(v)
					o.trigger.Update(v)
				}
			case <-control.C:
				control.D <- true
				return
			}
		}
	}()
}
