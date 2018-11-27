package filters

// Filter defines the interface that
// filters or processes incoming data
// and decides when to notify the observers.
type Filter interface {
	//Check should return true if the observers should be notified.
	Check(interface{}) bool
	//Update processes the new value. Its return value is sent to the observers
	//and it is only called when Check returns true.
	Update(interface{}) interface{}
}

// None struct implements the Filter interface.
// It forwards all data unfiltered and unprocessed.
type None struct {
}

//Check always returns true.
func (t *None) Check(newValue interface{}) bool { return true }

//Update returns the current value that is sent to the observers.
func (t *None) Update(newValue interface{}) interface{} {
	return newValue
}

// OnChange struct implements the Filter interface.
// It triggers when the value under observation changes.
type OnChange struct {
	Value interface{}
}

//Check if new value is different from saved value.
func (t *OnChange) Check(newValue interface{}) bool { return t.Value != newValue }

//Update updates stored value with new value and sends it to the observers.
func (t *OnChange) Update(newValue interface{}) interface{} {
	t.Value = newValue
	return newValue
}

// OnValue struct implements the Filter interface.
// It triggers when a certain value is reached.
type OnValue struct {
	Value interface{}
}

//Check returns true when stored value matches incoming value.
func (t *OnValue) Check(newValue interface{}) bool { return t.Value == newValue }

//Update returns the incoming value.
func (t *OnValue) Update(newValue interface{}) interface{} {
	return newValue
}

// AboveFloat64 struct implements the Filter interface.
// It triggers when a value is above a predefined value.
type AboveFloat64 struct {
	Value float64
}

//Check compares the stored to the incoming value and returns true if new value is above the stored one.
func (t *AboveFloat64) Check(newValue interface{}) bool { return newValue.(float64) > t.Value }

//Update returns incoming value.
func (t *AboveFloat64) Update(newValue interface{}) interface{} {
	return newValue
}

// BelowFloat64 struct implements the Filter interface.
// It triggers when a value is below a predefined value.
type BelowFloat64 struct {
	Value float64
}

//Check compares the stored to the incoming value and returns true if new value is below the stored one.
func (t *BelowFloat64) Check(newValue interface{}) bool { return newValue.(float64) < t.Value }

//Update returns incoming value.
func (t *BelowFloat64) Update(newValue interface{}) interface{} {
	return newValue
}

// MovingAverage implements the Filter interface.
// It requires a Size parameter to be initialized.
type MovingAverage struct {
	Size   int
	values []float64
}

//NewMovingAverage returns a moving average filter.
func NewMovingAverage(size int) Filter {
	return &MovingAverage{Size: size}
}

//Check returns always true because every value needs to be processed.
func (t *MovingAverage) Check(newValue interface{}) bool { return true }

//Update stores the new value and returns the moving average of updated data set.
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
