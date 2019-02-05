package filters

import (
	"fmt"
	"io"
	"math"
	"time"
)

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

// Model struct implements the Filter interface.
// It forwards all data unfiltered and unprocessed.
// Model can be embedded in structs to write user-defined filters.
type Model struct {
}

//Check always returns true.
func (t *Model) Check(newValue interface{}) bool { return true }

//Update returns the current value that is sent to the observers.
func (t *Model) Update(newValue interface{}) interface{} {
	return newValue
}

// None struct implements the Filter interface.
// It forwards all data unfiltered and unprocessed.
type None struct {
	Model
}

// Sink never forwards any values;
// it essentially blocks the flow of data
type Sink struct {
	Model
}

//Check always returns false
func (s *Sink) Check(newValue interface{}) bool {
	return false
}

// Mute struct implements the Filter interface.
// It blocks the forwarding of values within a
// predefined time duration.
type Mute struct {
	Period   time.Duration
	previous time.Time
	Model
}

//Check always returns true.
func (m *Mute) Check(newValue interface{}) bool {
	t := time.Now()
	if t.Sub(m.previous) < m.Period {
		return false
	}
	m.previous = t
	return true
}

// Print struct implements the Filter interface.
// It writes the incoming values to an io.Writer with a prefix.
type Print struct {
	Writer io.Writer
	Prefix string
	Model
}

//Update forwards value and prints it to io.Writer.
func (p *Print) Update(v interface{}) interface{} {
	_, err := fmt.Fprintf(p.Writer, "%s %v\n", p.Prefix, v)
	if err != nil {
		fmt.Printf("%s %v\n", p.Prefix, v)
	}
	return v
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
	Window int
	values []float64
}

//NewMovingAverage returns a moving average filter.
func NewMovingAverage(size int) Filter {
	return &MovingAverage{Window: size}
}

//Check returns always true because every value needs to be processed.
func (t *MovingAverage) Check(newValue interface{}) bool { return true }

//Update stores the new value and returns the moving average of updated data set.
func (t *MovingAverage) Update(newValue interface{}) interface{} {
	t.values = append(t.values, newValue.(float64))
	if len(t.values) > t.Window {
		t.values = t.values[len(t.values)-t.Window:]
	}
	var movingAverage float64
	for _, v := range t.values {
		movingAverage += v
	}
	return movingAverage / float64(len(t.values))
}

// Sigma implements the Filter interface.
// Check calculates the number of standard deviations that
// the incoming value is away from the mean.
// If that value is above the defined factor, it notifies the subscribers.
// Update returns the incoming value.
type Sigma struct {
	Window int
	Factor float64
	values []float64
	stddev float64
	mean   float64
	Model
}

//Check returns always true because every value needs to be processed.
func (s *Sigma) Check(newValue interface{}) bool {
	value := newValue.(float64)
	if len(s.values) >= s.Window {
		if s.stddev > 0.0 {
			sigma := (value - s.mean) / s.stddev
			if math.Abs(sigma) > s.Factor {
				return true
			}
		}
	}
	s.values = append(s.values, value)
	if len(s.values) > s.Window {
		s.values = s.values[len(s.values)-s.Window:]
	}
	var sum, sum2 float64
	for _, v := range s.values {
		sum += v
		sum2 += v * v
	}
	total := float64(len(s.values))
	s.mean = sum / total
	s.stddev = math.Sqrt((total*sum2 - sum*sum) / (total * total)) // sqrt(E[x*x] - E[x]^2)
	//fmt.Printf("mean=%+3.3f sd=%+3.3f\n", s.mean, s.stddev)
	return false
}

// Stddev implements the Filter interface.
// Check returns always true for data processing
// Update calculates standard deviation for a given window
// for every incoming data point and returns this value.
type Stddev struct {
	Window int
	values []float64
	stddev float64
	Model
}

//Update returns the standard deviation of samples in the window.
func (s *Stddev) Update(newValue interface{}) interface{} {
	value := newValue.(float64)
	s.values = append(s.values, value)
	if len(s.values) > s.Window {
		s.values = s.values[len(s.values)-s.Window:]
	}
	var sum, sum2 float64
	for _, v := range s.values {
		sum += v
		sum2 += v * v
	}
	total := float64(len(s.values))
	return math.Sqrt((total*sum2 - sum*sum) / (total * total)) // sqrt(E[x*x] - E[x]^2)
}

// Lowpass implements exponential smoothing.
// Can be used as RC lowpass filter when a = dt / (RC + dt).
// Check returns always true for data processing
// Update returns the filtered values.
type Lowpass struct {
	A        float64
	oldValue float64
	Model
}

//Update returns the exponentially smoothened parameter
func (lp *Lowpass) Update(newValue interface{}) interface{} {
	x := newValue.(float64)
	lpValue := lp.A*x + (1-lp.A)*lp.oldValue
	lp.oldValue = lpValue
	return lpValue
}
