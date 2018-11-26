package notifiers

// Notifier defines the interface
// for a trigger that is required
// to create a new observer.
type Notifier interface {
	Check(interface{}) bool
	Update(interface{})
}

// OnChange struct implements the Notifier interface.
// It triggers when the value under observation changes.
type OnChange struct {
	Value interface{}
}

//Check returns true if the observers should be notified
func (t *OnChange) Check(newValue interface{}) bool { return t.Value != newValue }

//Update handles a new value depending on the trigger implementation
func (t *OnChange) Update(newValue interface{}) { t.Value = newValue }

// OnValue struct implements the Notifier interface.
// It triggers when a certain value is reached.
type OnValue struct {
	Value interface{}
}

//Check returns true if the observers should be notified
func (t *OnValue) Check(newValue interface{}) bool { return t.Value == newValue }

//Update handles a new value depending on the trigger implementation
func (t *OnValue) Update(newValue interface{}) {}

// AboveFloat64 struct implements the Notifier interface.
// It triggers when a value is above a predefined value
type AboveFloat64 struct {
	Value float64
}

//Check returns true if the observers should be notified
func (t *AboveFloat64) Check(newValue interface{}) bool { return newValue.(float64) > t.Value }

//Update handles a new value depending on the trigger implementation
func (t *AboveFloat64) Update(newValue interface{}) {}

// BelowFloat64 struct implements the trigger interface.
// It triggers when a value is below a predefined value
type BelowFloat64 struct {
	Value float64
}

//Check returns true if the observers should be notified
func (t *BelowFloat64) Check(newValue interface{}) bool { return newValue.(float64) < t.Value }

//Update handles a new value depending on the trigger implementation
func (t *BelowFloat64) Update(newValue interface{}) {}
