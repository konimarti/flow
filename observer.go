package observer

// observer interface
type Observer interface {
	Notify(interface{})
	Subscribe() chan interface{}
	Unsubscribe(chan interface{})
	Close()
}

//observerImpl implements the observer interface.
//This should be used in type embedding.
type observerImpl struct {
	trigger   Trigger
	observers []chan interface{}
	closing   []*control
}

//Notify sends out the current value in the observer channel
func (o *observerImpl) Notify(value interface{}) {
	for _, observer := range o.observers {
		select {
		case v := <-observer:
			// in case channel has been closed,
			// remove it from the list
			if v == nil {
				o.Unsubscribe(observer)
				continue
			}
		default:
		}
		observer <- value
	}
}

func (o *observerImpl) Subscribe() chan interface{} {
	observer := make(chan interface{}, 1)
	o.observers = append(o.observers, observer)
	return observer
}

func (o *observerImpl) Unsubscribe(ch chan interface{}) {
	// not sure if needed
}

//Close closes all the observers channels
func (o *observerImpl) Close() {
	for _, control := range o.closing {
		control.Close()
	}
}
