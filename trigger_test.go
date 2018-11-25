package observer_test

import (
	"testing"

	"github.com/konimarti/observer"
)

var configT = []struct {
	Value  interface{}
	Update interface{}
}{
	{
		Value:  1.0,
		Update: 2.0,
	},
	{
		Value:  1,
		Update: 2,
	},
	{
		Value:  "hello",
		Update: "world",
	},
}

func TestOnChange(t *testing.T) {
	for _, cfg := range configT {
		//new trigger
		trig := observer.OnChange{}
		trig.Update(cfg.Value)

		//test fire
		if trig.Fire(cfg.Value) {
			t.Error("should not fire because value are the same.")
		}
		if !trig.Fire(cfg.Update) {
			t.Error("should not fire because value are not the same.")
		}

		//test update
		trig.Update(cfg.Update)
		if trig.Fire(cfg.Update) {
			t.Error("should not fire because value are the same after update.")
		}
	}
}

func TestOnValue(t *testing.T) {
	for _, cfg := range configT {
		//new trigger
		trig := observer.OnValue{Value: cfg.Value}

		//test fire
		if !trig.Fire(cfg.Value) {
			t.Error("should fire because value are the same.")
		}
		if trig.Fire(cfg.Update) {
			t.Error("should not fire because value are not the same.")
		}

		//test update
		trig.Update(cfg.Update)
		if trig.Fire(cfg.Update) {
			t.Error("should not fire because update does not change the value.")
		}
	}
}
