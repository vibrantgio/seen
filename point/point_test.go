package point

import (
	"math"
	"testing"

	"github.com/vibrantgio/seen/float"
	"github.com/vibrantgio/seen/quaternion"
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

func TestAxisAngle(t *testing.T) {
	q1 := Pt(2, 0, 0).PointAngle(math.Pi / 2.0)
	if !float.Equal(q1.Length(), 1.0) {
		t.Errorf("Exp: q1.Length()==1.0\nGot: q1.Length()==%v", q1.Length())
	}

	v := quaternion.Q(0, 1, 0, 1)
	v = q1.Mul(v).Mul(q1.Conjugate())

	if !float.EqualPairs(v.X, 0, v.Y, 0, v.Z, 1, v.W, 1) {
		t.Errorf("Exp: {0,0,1,1}\nGot: %v", v)
	}

	vx, vy, vz := q1.Transform(0, 1, 0)
	if !float.EqualPairs(vx, 0, vy, 0, vz, 1) {
		t.Errorf("Exp: {0,0,1}\nGot: {%v,%v,%v}", vx, vy, vz)
	}
}
