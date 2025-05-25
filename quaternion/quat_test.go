package quaternion_test

import (
	"math"
	"testing"

	"github.com/vibrantgio/seen/float"
	"github.com/vibrantgio/seen/quaternion"
)

func TestInitialization(t *testing.T) {

	q1 := quaternion.Identity

	if !float.Equal(q1.W, 1) {
		t.Fail()
	}

	q2 := quaternion.Quat{0, 0, 0, 1}

	if !q1.Equal(q2) {
		t.Fail()
	}
}

func TestEqual(t *testing.T) {

	q1 := quaternion.Q(1, 2, 3, 4)
	q2 := quaternion.Q(5, 6, 7, 8)

	if !q1.Equal(q1) {
		t.Fail()
	}

	if !q2.Equal(q2) {
		t.Fail()
	}

	if q1.Equal(q2) {
		t.Fail()
	}
}

func TestConjugate(t *testing.T) {

	q := quaternion.Q(1, 2, 3, 4)
	c := q.Conjugate()

	if !float.Equal(q.X, -c.X) {
		t.Fail()
	}

	if !float.Equal(q.Y, -c.Y) {
		t.Fail()
	}
	if !float.Equal(q.Z, -c.Z) {
		t.Fail()
	}
	if !float.Equal(q.W, c.W) {
		t.Fail()
	}
}

func TestAdd(t *testing.T) {

	q1 := quaternion.Q(1, 2, 3, 4)
	q2 := quaternion.Q(5, 6, 7, 8)
	expect := quaternion.Q(6, 8, 10, 12)

	if !q1.Add(q2).Equal(expect) {
		t.Fail()
	}
}

func TestScale(t *testing.T) {

	q := quaternion.Q(1, 2, 3, 4)
	q = q.Normalize()
	q = q.Scale(2)

	if !float.Equal(q.Length(), 2) {
		t.Fail()
	}
}

func TestDot(t *testing.T) {

	q1 := quaternion.Q(1, 2, 3, 4)
	d1 := q1.Dot(q1)

	if !float.Equal(d1, 30) {
		t.Fail()
	}

	q2 := quaternion.Q(5, 6, 7, 8)
	d2 := q1.Dot(q2)

	if !float.Equal(d2, 70) {
		t.Fail()
	}
}

func TestLength(t *testing.T) {

	q := quaternion.Q(1, 2, 3, 4)
	l := q.Length()

	if !float.Equal(l, math.Sqrt(q.Dot(q))) {
		t.Fail()
	}
}

func TestMul(t *testing.T) {

	// Multiplying a quaternion with its conjugate should leave only a real component.
	// So only W should hold a value and X,Y and Z should be zero.
	// The W value should be equal to the dot product of the components
	q := quaternion.Q(1, 2, 3, 4)
	prod := q.Mul(q.Conjugate())
	expect := quaternion.Q(0, 0, 0, q.Dot(q))

	if !prod.Equal(expect) {
		t.Fail()
	}
}

func TestNormalize(t *testing.T) {

	q := quaternion.Q(0, 0, 0, 2)
	r := q.Normalize()
	if !float.Equal(r.Dot(r), 1) {
		t.Fail()
	}

	q = quaternion.Q(1, 2, 3, 4)
	r = q.Normalize()
	if !float.Equal(r.Dot(r), 1) {
		t.Fail()
	}
}

func TestAxisAngle(t *testing.T) {
	q := quaternion.AxisAngle(2, 0, 0, math.Pi/2.0)
	if !float.Equal(q.Length(), 1.0) {
		t.Errorf("Exp: q.Length()==1.0\nGot: q.Length()==%v", q.Length())
	}

	v := quaternion.Q(0, 1, 0, 1)
	v = q.Mul(v).Mul(q.Conjugate())

	if !float.EqualPairs(v.X, 0, v.Y, 0, v.Z, 1, v.W, 1) {
		t.Errorf("Exp: {0,0,1,1}\nGot: %v", v)
	}

	vx, vy, vz := q.Transform(0, 1, 0)
	if !float.EqualPairs(vx, 0, vy, 0, vz, 1) {
		t.Errorf("Exp: {0,0,1}\nGot: {%v,%v,%v}", vx, vy, vz)
	}
}

func TestEvaluationOrer(t *testing.T) {

	rx := quaternion.RotX(0.1 * math.Pi)
	ry := quaternion.RotY(0.3 * math.Pi)
	rz := quaternion.RotZ(0.5 * math.Pi)

	ux, uy, uz := rz.Transform(1, 2, 3)
	ux, uy, uz = ry.Transform(ux, uy, uz)
	ux, uy, uz = rx.Transform(ux, uy, uz)

	vx, vy, vz := rx.Mul(ry).Mul(rz).Transform(1, 2, 3)
	wx, wy, wz := quaternion.RotX(0.1*math.Pi).RotY(0.3*math.Pi).RotZ(0.5*math.Pi).Transform(1, 2, 3)

	if !float.EqualPairs(ux, vx, uy, vy, uz, vz) {
		t.Errorf("Exp: {%v,%v,%v}\nGot: {%v,%v,%v}", ux, uy, uz, vx, vy, vz)
	}
	if !float.EqualPairs(ux, wx, uy, wy, uz, wz) {
		t.Errorf("Exp: {%v,%v,%v}\nGot: {%v,%v,%v}", ux, uy, uz, wx, wy, wz)
	}
}
