package filters_test

import (
	"testing"

	"github.com/konimarti/flow/filters"
)

func TestGetFloat64(t *testing.T) {
	input := []interface{}{
		int(1),
		int16(2),
		int32(3),
		int64(4),
		float32(5),
		float64(6),
		string("1"),
		struct{}{},
	}
	want := []float64{
		1.0,
		2.0,
		3.0,
		4.0,
		5.0,
		6.0,
		0.0,
		0.0,
	}

	for i, in := range input {
		if want[i] != filters.GetFloat64(in) {
			t.Error("GetFloat64 failed")
		}
	}
}
