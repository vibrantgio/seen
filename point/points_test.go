package point

import (
	"testing"

	"github.com/vibrantgio/seen/float"
)

func TestPointsRoundAlterInSlice(t *testing.T) {
	points := Points{{1.4, 1.5, 1.6}, {1.2, 1.5, 1.6}}

	// Calling a method that alters the object will alter the points in the slice.
	points[0] = points[0].Round()
	points[1] = points[1].Round()

	// Check that the points where indeed altered
	p := points[0]
	if !float.EqualPairs(p.X, 1, p.Y, 2, p.Z, 2) {
		t.Fail()
	}
	p = points[1]
	if !float.EqualPairs(p.X, 1, p.Y, 2, p.Z, 2) {
		t.Fail()
	}
}
func TestPointsRoundAlteredInFor(t *testing.T) {
	points := Points{{1.4, 1.5, 1.6}, {1.2, 1.5, 1.6}}

	// addressing the slice via index will alter the point in the slice.
	// So assigning to the point will alter it.
	for i, p := range points {
		points[i] = p.Round()
	}

	// Verify that points were not actually altered
	p := points[0]
	if !float.EqualPairs(p.X, 1, p.Y, 2, p.Z, 2) {
		t.Fail()
	}
	p = points[1]
	if !float.EqualPairs(p.X, 1, p.Y, 2, p.Z, 2) {
		t.Fail()
	}
}
