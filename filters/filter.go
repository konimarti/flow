package filters

// Filter defines the interface
// that filters incoming data
// and decides when to notify the observers.
type Filter interface {
	Check(interface{}) bool
	Update(interface{}) interface{}
}

// OnChange struct implements the Filter interface.
// It triggers when the value under observation changes.
type OnChange struct {
	Value interface{}
}

//Check returns true if the observers should be notified.
func (t *OnChange) Check(newValue interface{}) bool { return t.Value != newValue }

//Update handles a new value depending on the Filter.
func (t *OnChange) Update(newValue interface{}) interface{} {
	t.Value = newValue
	return newValue
}

// OnValue struct implements the Filter interface.
// It triggers when a certain value is reached.
type OnValue struct {
	Value interface{}
}

//Check returns true if the observers should be notified
func (t *OnValue) Check(newValue interface{}) bool { return t.Value == newValue }

//Update handles a new value depending on the Filter.
func (t *OnValue) Update(newValue interface{}) interface{} {
	return newValue
}

// AboveFloat64 struct implements the Filter interface.
// It triggers when a value is above a predefined value.
type AboveFloat64 struct {
	Value float64
}

//Check returns true if the observers should be notified.
func (t *AboveFloat64) Check(newValue interface{}) bool { return newValue.(float64) > t.Value }

//Update handles a new value depending on the Filter.
func (t *AboveFloat64) Update(newValue interface{}) interface{} {
	return newValue
}

// BelowFloat64 struct implements the Filter interface.
// It triggers when a value is below a predefined value.
type BelowFloat64 struct {
	Value float64
}

//Check returns true if the observers should be notified.
func (t *BelowFloat64) Check(newValue interface{}) bool { return newValue.(float64) < t.Value }

//Update handles a new value depending on the Filter.
func (t *BelowFloat64) Update(newValue interface{}) interface{} {
	return newValue
}

// MovingAverage struct implements the Filter interface.
// It triggers when a value is below a predefined value.
type MovingAverage struct {
	Size   int
	values []float64
}

//Check returns true if the observers should be notified.
func (t *MovingAverage) Check(newValue interface{}) bool { return true }

//Update handles a new value depending on the Filter.
func (t *MovingAverage) Update(newValue interface{}) interface{} {
	t.values = append(t.values, newValue.(float64))
	if len(t.values) > t.Size {
		t.values = t.values[len(t.values)-t.Size:]
	}
	var movingAverage float64
	for _, v := range t.values {
		movingAverage += v
	}
	return movingAverage / float64(len(t.values))
}
