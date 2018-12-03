package filters_test

import (
	"fmt"
	"testing"

	"github.com/konimarti/observer/filters"
)

func TestSwitch(t *testing.T) {
	var config = []struct {
		Name   string
		Values []interface{}
		Left   filters.Filter
		Right  filters.Filter
		Checks []bool
		Wants  []interface{}
	}{
		{
			Name:   "Float64Collar",
			Values: []interface{}{0.0, 0.5, 1.0, 1.1, -1.0, -1.2, 0.0},
			Left:   &filters.AboveFloat64{1.0},
			Right:  &filters.BelowFloat64{-1.0},
			Checks: []bool{false, false, false, true, false, true, false},
			Wants:  []interface{}{nil, nil, nil, 1.1, nil, -1.2, nil},
		},
	}

	for _, cfg := range config {
		// internal consistency checks
		or := filters.NewSwitch(cfg.Left, cfg.Right)
		if len(cfg.Values) != len(cfg.Checks) {
			t.Error("internal test consistency failed (values vs. Checks)")
		}
		if len(cfg.Values) != len(cfg.Wants) {
			t.Error("internal test consistency failed (values vs. Wants)")
		}
		// process values stream and check results for check and update
		for i, _ := range cfg.Values {
			check := or.Check(cfg.Values[i])
			if check != cfg.Checks[i] {
				fmt.Printf("Name: %s. Got %v. Expected %v.\n", cfg.Name, check, cfg.Checks[i])
				t.Error("check failed")
			}
			if check {
				value := or.Update(cfg.Values[i])
				if value != cfg.Wants[i] {
					fmt.Printf("Name: %s. Got %v. Expected %v.\n", cfg.Name, value, cfg.Wants[i])
					t.Error("update failed")
				}
			}
		}
	}
}
