package filters_test

import (
	"testing"

	"github.com/konimarti/observer/filters"
)

func TestChain(t *testing.T) {
	var config = []struct {
		Values []interface{}
		Chain  []filters.Filter
		Checks []bool
		Wants  []interface{}
	}{
		{
			Values: []interface{}{1, 1, 2, 2, 2, 3, 3, 3, 3, 4},
			Chain:  []filters.Filter{&filters.OnChange{}, &filters.OnValue{3}},
			Checks: []bool{false, false, false, false, false, true, false, false, false, false},
			Wants:  []interface{}{nil, nil, nil, nil, nil, 3, nil, nil, nil, nil},
		},
		{
			Values: []interface{}{1.1, 1.1, 2.1, 2.1, 2.1, 3.5, 3.5, 3.5, 3.5, 4.0},
			Chain:  []filters.Filter{&filters.OnChange{}, &filters.OnValue{3.0}},
			Checks: []bool{false, false, false, false, false, true, false, false, false, false},
			Wants:  []interface{}{nil, nil, nil, nil, nil, 3.5, nil, nil, nil, nil},
		},
		{
			Values: []interface{}{"hello", "hello", "world", "world"},
			Chain:  []filters.Filter{&filters.OnChange{}, &filters.OnValue{"world"}},
			Checks: []bool{false, false, true, false},
			Wants:  []interface{}{nil, nil, "world", nil},
		},
		{
			Values: []interface{}{1.0, 2.0, 2.0, 3.5, 5.0, 6.0},
			Chain:  []filters.Filter{&filters.None{}, &filters.AboveFloat64{3.0}, &filters.BelowFloat64{4.0}},
			Checks: []bool{false, false, false, true, false, false},
			Wants:  []interface{}{nil, nil, nil, 3.5, nil, nil},
		},
	}

	for _, cfg := range config {
		// internal consistency checks
		chain := filters.NewChain(cfg.Chain...)
		if len(cfg.Values) != len(cfg.Checks) {
			t.Error("internal test consistency failed (values vs. Checks)")
		}
		if len(cfg.Values) != len(cfg.Wants) {
			t.Error("internal test consistency failed (Values vs. Wants)")
		}
		// process values stream and check results for check and update
		for i, _ := range cfg.Values {
			check := chain.Check(cfg.Values[i])
			if check != cfg.Checks[i] {
				t.Error("check failed. Got %v. Expected %v.", check, cfg.Checks[i])
			}
			if check {
				value := chain.Update(cfg.Values[i])
				if value != cfg.Wants[i] {
					t.Error("update failed. Got %v. Expected %v.", value, cfg.Wants[i])
				}
			}
		}
	}
}
