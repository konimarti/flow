package observer

type observerChannel struct {
	observerImpl
}

//NewFromChannel creates a new observer struct
func NewFromChannel(tr Trigger, channel chan interface{}) Observer {
	obs := observerChannel{
		observerImpl: observerImpl{
			control: NewControl(),
			trigger: tr,
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
				if o.observerImpl.trigger.Check(v) {
					o.Notify(v)
					o.trigger.Update(v)
				}
			case <-o.control.C:
				o.control.D <- true
				return
			}
		}
	}()
}
