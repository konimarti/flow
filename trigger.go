package observer

type Trigger interface {
	Fire(interface{}) bool
	Update(interface{})
}

// OnChange struct implements the trigger interface.
// It triggers when the value under observation changes.
type OnChange struct {
	Value interface{}
}

func (t *OnChange) Fire(newValue interface{}) bool { return t.Value != newValue }
func (t *OnChange) Update(newValue interface{})    { t.Value = newValue }

// OnValue struct implements the trigger interface.
// It triggers when a certain value is reached.
type OnValue struct {
	Value interface{}
}

func (t *OnValue) Fire(newValue interface{}) bool { return t.Value == newValue }
func (t *OnValue) Update(newValue interface{})    {}

// AboveFloat64 struct implements the trigger interface.
// It triggers when a value is above a predefined value
type AboveFloat64 struct {
	Value float64
}

func (t *AboveFloat64) Fire(newValue interface{}) bool { return newValue.(float64) > t.Value }
func (t *AboveFloat64) Update(newValue interface{})    {}
