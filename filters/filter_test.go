package filters_test

import (
	"bytes"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/konimarti/flow/filters"
)

var configT = []struct {
	Type   string
	Value  interface{}
	Update interface{}
}{
	{
		Type:   "float64",
		Value:  1.0,
		Update: 2.0,
	},
	{
		Type:   "int",
		Value:  1,
		Update: 2,
	},
	{
		Type:   "string",
		Value:  "hello",
		Update: "world",
	},
}

func TestNone(t *testing.T) {
	for _, cfg := range configT {
		//new trigger
		trig := filters.None{}
		trig.Update(cfg.Value)

		//test check
		if !trig.Check(cfg.Value) {
			t.Error("should fire because all values are processed")
		}
		if !trig.Check(cfg.Update) {
			t.Error("should fire because all values are processed")
		}

		//test update
		if cfg.Update != trig.Update(cfg.Update) && cfg.Value != trig.Update(cfg.Value) {
			t.Error("should return same values as it was called with")
		}
	}
}

func TestSink(t *testing.T) {
	for _, cfg := range configT {
		//new trigger
		trig := filters.Sink{}

		//test check
		if trig.Check(cfg.Value) {
			t.Error("should not fire because no values are processed")
		}
	}
}

func TestPrint(t *testing.T) {
	for _, cfg := range configT {
		//new trigger
		var buf bytes.Buffer
		prefix := "Prefix"
		trig := filters.Print{Writer: &buf, Prefix: prefix}

		//test check
		if !trig.Check(cfg.Value) {
			t.Error("should fire because all values are processed")
		}
		if !trig.Check(cfg.Update) {
			t.Error("should fire because all values are processed")
		}

		//test update
		if cfg.Value != trig.Update(cfg.Value) {
			t.Error("should return same values as it was called with")
		}
		str := fmt.Sprintf("%s %v\n", prefix, cfg.Value)
		if buf.String() != str {
			fmt.Printf("Got: %s. Expected: %s\n", buf.String(), str)
			t.Error("Printed string does not match")
		}
	}
}

func TestOnChange(t *testing.T) {
	for _, cfg := range configT {
		//new trigger
		trig := filters.OnChange{}
		trig.Update(cfg.Value)

		//test check
		if trig.Check(cfg.Value) {
			t.Error("should not fire because value are the same.")
		}
		if !trig.Check(cfg.Update) {
			t.Error("should not fire because value are not the same.")
		}

		//test update
		trig.Update(cfg.Update)
		if trig.Check(cfg.Update) {
			t.Error("should not fire because value are the same after update.")
		}
	}
}

func TestOnValue(t *testing.T) {
	for _, cfg := range configT {
		//new trigger
		trig := filters.OnValue{Value: cfg.Value}

		//test check
		if !trig.Check(cfg.Value) {
			t.Error("should fire because value are the same.")
		}
		if trig.Check(cfg.Update) {
			t.Error("should not fire because value are not the same.")
		}

		//test update
		trig.Update(cfg.Update)
		if trig.Check(cfg.Update) {
			t.Error("should not fire because update does not change the value.")
		}
	}
}

func TestAboveFloat64(t *testing.T) {
	for _, cfg := range configT {
		if cfg.Type != "float64" {
			continue
		}
		//new trigger
		trig := filters.AboveFloat64{Value: cfg.Value.(float64)}

		//test check
		if trig.Check(cfg.Value) {
			t.Error("should not fire because value are the same.")
		}
		if !trig.Check(cfg.Update) {
			t.Error("should fire because new value is above old value.")
		}

		//test update
		trig.Update(cfg.Update)
		if trig.Value != cfg.Value {
			t.Error("update should not change initial, predefined value.")
		}
	}
}

func TestBelowFloat64(t *testing.T) {
	for _, cfg := range configT {
		if cfg.Type != "float64" {
			continue
		}
		//new trigger
		trig := filters.BelowFloat64{Value: cfg.Update.(float64)}

		//test check
		if trig.Check(cfg.Update) {
			t.Error("should not fire because value are the same.")
		}
		if !trig.Check(cfg.Value) {
			t.Error("should fire because new value is below old value.")
		}

		//test update
		trig.Update(cfg.Value)
		if trig.Value != cfg.Update {
			t.Error("update should not change initial, predefined value.")
		}
	}
}

func TestMovingAverage(t *testing.T) {
	values := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0}
	mvavg := []float64{1.0, 1.5, 2.0, 2.5, 3.0, 4.0, 5.0}

	//new trigger
	trig := filters.NewMovingAverage(5)

	for i, v := range values {

		//test check
		if !trig.Check(v) {
			t.Error("should always fire")
		}

		//test update
		mv := trig.Update(v)
		if math.Abs(mvavg[i]-mv.(float64)) > 1e-6 {
			t.Error("moving average not calculated correctly")
		}
	}
}

func TestSigma(t *testing.T) {
	values := []float64{1.0, 1.1, 2.0, 1.05}
	checks := []bool{false, false, true, false}

	//new trigger
	trig := filters.Sigma{Window: 2, Factor: 1.0}

	for i, v := range values {

		//test check
		c := trig.Check(v)
		if checks[i] != c {
			fmt.Printf("Got %v. Expected %v", c, checks[i])
			t.Error("check failed")
		}

		//test update
		value := trig.Update(v)
		if math.Abs(v-value.(float64)) > 1e-6 {
			t.Error("updated value should be the same")
		}
	}
}

func TestStddev(t *testing.T) {
	values := []float64{1.0, 1.1, 2.0, 1.05}
	checks := []bool{true, true, true, true}
	processed := []float64{0.0, 0.05, 0.4496913, 0.4365267}

	//new trigger
	trig := filters.Stddev{Window: 3}

	for i, v := range values {

		//test check
		c := trig.Check(v)
		if checks[i] != c {
			fmt.Printf("Got %v. Expected %v", c, checks[i])
			t.Error("check failed")
		}

		//test update
		value := trig.Update(v)
		if math.Abs(processed[i]-value.(float64)) > 1e-6 {
			fmt.Printf("Got %v. Expected %v", value.(float64), processed[i])
			t.Error("updated value should be the same")
		}
	}
}

func TestMute(t *testing.T) {
	values := []interface{}{1.0, 1, "hello"}

	for _, v := range values {
		//new trigger
		trig := filters.Mute{Duration: 100 * time.Millisecond}

		b1 := trig.Check(v)
		b2 := trig.Check(v)
		time.Sleep(100 * time.Millisecond)
		b3 := trig.Check(v)
		if b1 != true {
			t.Error("first check should be true")
		}
		if b2 != false {
			t.Error("second check immediately after first should be false")
		}
		if b3 != true {
			t.Error("third check after mute period should be true")
		}

		// Update
		newVal := trig.Update(v)
		if v != newVal {
			fmt.Printf("Got %v. Expected %v", newVal, v)
			t.Error("value returend from update should be the same")
		}
	}
}

func TestLowPass(t *testing.T) {
	a := []float64{0.0, 1.0, 0.5}
	input := []float64{1.0, 2.0, 3.0}
	expected := [][]float64{
		[]float64{0.0, 0.0, 0.0},
		[]float64{1.0, 2.0, 3.0},
		[]float64{0.5, 1.25, 2.125},
	}

	for i, factor := range a {
		trig := filters.LowPass{A: factor}
		for j, in := range input {
			//test check
			if !trig.Check(in) {
				t.Error("should always fire")
			}
			// test update
			out := trig.Update(in)
			if math.Abs(expected[i][j]-out.(float64)) > 1e-6 {
				fmt.Println("failed for i=", i, "j=", j)
				fmt.Println("input =", in)
				fmt.Println("expected =", expected[i][j])
				fmt.Println("received =", out.(float64))
				t.Error("lowpass filter does not work correctly")
			}
		}
	}
}
