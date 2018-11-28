package filters_test

import (
	"bytes"
	"fmt"
	"math"
	"testing"

	"github.com/konimarti/observer/filters"
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
