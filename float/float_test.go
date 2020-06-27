package float

import (
	"math"
	"testing"
)

func TestFloatThreshold(t *testing.T) {

	var a, b, epsilon float64 = 0.00001 + 2.2e-16, 0.00001, Epsilon

	t.Log("diff (absolute error):", math.Abs(a-b), "diff/Abs(b) (relative error):", math.Abs(a-b)/math.Abs(b), "epsilon:", epsilon)

	diff := math.Abs(a - b)
	if diff/math.Min(math.Abs(b), MaxValue) >= epsilon {
		t.Fail()
	}
}
