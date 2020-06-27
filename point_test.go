package seen

import (
	"testing"
	"github.com/reactivego/seen/float"
)

func TestPointRoundReturnedPointIsNewInstance(t *testing.T) {
	p := &Point{1.1, 2.1, 3.1} // original pointer to a point
	q := p.Round()                // should be a unique instance
	z := p.Round()                // should again be a unique instance

	// make sure the pointers are unique
	if q == p || q == z {
		t.Fail()
	}

	// Make sure the rounded value is different from the original value
	if *p == *q {
		t.Fail()
	}

	// Make sure the rounded values are identical
	if *q != *z {
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
	points[0].RoundAssign()
	points[1].RoundAssign()

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

func TestPointRoundNotAlteredInFor(t *testing.T) {
	points := []Point{{1.4, 1.5, 1.6}, {1.2, 1.5, 1.6}}

	// p is expected to be a copy of the point in the slice.
	// So modifying p will not alter it in the slice.
	for _, p := range points {
		p.RoundAssign()
	}

	// Verify that points were not actually altered
	p := points[0]
	if !float.EqualPairs(p.X, 1.4, p.Y, 1.5, p.Z, 1.6) {
		t.Fail()
	}
	p = points[1]
	if !float.EqualPairs(p.X, 1.2, p.Y, 1.5, p.Z, 1.6) {
		t.Fail()
	}
}

func TestPointRoundAlteredInFor(t *testing.T) {
	points := []Point{{1.4, 1.5, 1.6}, {1.2, 1.5, 1.6}}

	// addressing the slice via index will alter the point in the slice.
	// So calling RoundAssign() on the point will alter it.
	for i := range points {
		points[i].RoundAssign()
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
	x := MakePointX()
	y := MakePointY()
	// Geometric interpretation of x cross y is to rotate x to y using right hand
	// rule, thumb points in the direction of resulting vector.
	// Cross product of two axis of orthogonal coordinate system should produce the third axis.
	c := x.Cross(y)
	t.Log(c)
	if !c.Equal(MakePointZ()) {
		t.Fail()
	}
}
