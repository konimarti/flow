package observer_test

import (
	"testing"

	"github.com/konimarti/observer"
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

func TestOnChange(t *testing.T) {
	for _, cfg := range configT {
		//new trigger
		trig := observer.OnChange{}
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
		trig := observer.OnValue{Value: cfg.Value}

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
		trig := observer.AboveFloat64{Value: cfg.Value.(float64)}

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
		trig := observer.BelowFloat64{Value: cfg.Update.(float64)}

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
