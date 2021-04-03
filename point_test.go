package seen

import (
	"testing"

	"github.com/reactivego/seen/float"
)

func TestPointRoundReturnedPointIsUniqueValue(t *testing.T) {
	p := Point{1.1, 2.1, 3.1} // original point value
	q := p.Round()            // should be a unique value
	z := p.Round()            // should again be a unique value

	// Make sure the rounded value is different from the original value
	if p == q {
		t.Fail()
	}

	// Make sure the rounded values are identical
	if q != z {
		t.Fail()
	}

	// Use the Equal method to verify rounded values are identical
	if !q.Equal(z) {
		t.Fail()
	}
}

func TestPointRoundAlterInSlice(t *testing.T) {
	points := []Point{{1.4, 1.5, 1.6}, {1.2, 1.5, 1.6}}

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
func TestPointRoundAlteredInFor(t *testing.T) {
	points := []Point{{1.4, 1.5, 1.6}, {1.2, 1.5, 1.6}}

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

func TestPointCross(t *testing.T) {
	x := PointX
	y := PointY
	// Geometric interpretation of x cross y is to rotate x to y using right hand
	// rule, thumb points in the direction of resulting vector.
	// Cross product of two axis of orthogonal coordinate system should produce the third axis.
	c := x.Cross(y)
	t.Log(c)
	if !c.Equal(PointZ) {
		t.Fail()
	}
}
