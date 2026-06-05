package svg

import (
	"strconv"
	"testing"
)

func TestFormatFloat(t *testing.T) {
	if strconv.FormatFloat(234.553343, 'g', -1, 64) != "234.553343" {
		t.Fail()
	}
}

func TestFtoaPrecision(t *testing.T) {
	cases := []struct {
		precision int
		v         float64
		want      string
	}{
		{0, 123.7, "124"},        // round to nearest integer (up)
		{0, 123.2, "123"},        // round to nearest integer (down)
		{2, 123.459, "123.46"},   // keep two decimals
		{3, 123.4567, "123.457"}, // keep three decimals
		{-1, 0.5, "0.5"},         // shortest round-tripping form
	}
	for _, c := range cases {
		if got := Ftoa(c.precision, c.v); got != c.want {
			t.Errorf("Ftoa(%d, %v) = %q, want %q", c.precision, c.v, got, c.want)
		}
	}
}
