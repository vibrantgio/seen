package seen

import (
	"testing"
)

func TestPointRoundReturnedPointIsUniqueValue(t *testing.T) {
	p := Point{X: 1.1, Y: 2.1, Z: 3.1} // original point value
	q := p.Round()                     // should be a unique value
	z := p.Round()                     // should again be a unique value

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

func TestPointCross(t *testing.T) {
	x := Point{X: 1}
	y := Point{Y: 1}
	// Geometric interpretation of x cross y is to rotate x to y using right hand
	// rule, thumb points in the direction of resulting vector.
	// Cross product of two axis of orthogonal coordinate system should produce the third axis.
	c := x.Cross(y)
	if !c.Equal(Point{Z: 1}) {
		t.Log(c)
		t.Fail()
	}
}
