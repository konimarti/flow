package observer_test

import "github.com/konimarti/observer"

var config = []struct {
	Values []interface{}
	Want   interface{}
}{
	{
		Values: []interface{}{1.1, 1.1, 1.1, 2.1, 1.1},
		Want:   2.1,
	},
	{
		Values: []interface{}{1, 1, 1, 2, 1},
		Want:   2,
	},
	{
		Values: []interface{}{"hello", "hello", "hello", "world", "hello"},
		Want:   "world",
	},
}

var observers = []struct {
	Name   string
	TrFunc func(v interface{}) observer.Trigger
}{
	{
		Name: "OnChange",
		TrFunc: func(v interface{}) observer.Trigger {
			return &observer.OnChange{v}
		},
	},
	{
		Name: "OnValue",
		TrFunc: func(v interface{}) observer.Trigger {
			return &observer.OnValue{v}
		},
	},
}
