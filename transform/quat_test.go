package transform

import (
	"math"
	"testing"

	"github.com/reactivego/seen/float"
)

func TestInitialization(t *testing.T) {
	q1 := IdentQuaternion

	if !float.Equal(q1.W, 1) {
		t.Fail()
	}

	q2 := Quaternion{0, 0, 0, 1}

	if !q1.Equal(q2) {
		t.Fail()
	}
}

func TestEqual(t *testing.T) {
	q1 := Quaternion{1, 2, 3, 4}
	q2 := Quaternion{5, 6, 7, 8}

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
	q := Quaternion{1, 2, 3, 4}
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
	q1 := Quaternion{1, 2, 3, 4}
	q2 := Quaternion{5, 6, 7, 8}
	expect := Quaternion{6, 8, 10, 12}
	if !q1.Add(q2).Equal(expect) {
		t.Fail()
	}

}

func TestScale(t *testing.T) {
	q := Quaternion{1, 2, 3, 4}
	q = q.Normalize()
	q = q.Scale(2)
	if !float.Equal(q.Length(), 2) {
		t.Fail()
	}
}

func TestDot(t *testing.T) {
	q1 := Quaternion{1, 2, 3, 4}
	d1 := q1.Dot(q1)

	if !float.Equal(d1, 30) {
		t.Fail()
	}

	q2 := Quaternion{5, 6, 7, 8}
	d2 := q1.Dot(q2)

	if !float.Equal(d2, 70) {
		t.Fail()
	}
}

func TestLength(t *testing.T) {
	q := Quaternion{1, 2, 3, 4}
	l := q.Length()

	if !float.Equal(l, math.Sqrt(q.Dot(q))) {
		t.Fail()
	}
}

func TestMul(t *testing.T) {
	// Multiplying a quaternion with its conjugate should leave only a real component.
	// So only W should hold a value and X,Y and Z should be zero.
	// The W value should be equal to the dot product of the components
	q := Quaternion{1, 2, 3, 4}
	prod := q.Mul(q.Conjugate())
	expect := Quaternion{0, 0, 0, q.Dot(q)}
	if !prod.Equal(expect) {
		t.Fail()
	}
}

func TestNormalize(t *testing.T) {
	q := Quaternion{0, 0, 0, 2}
	r := q.Normalize()
	if !float.Equal(r.Dot(r), 1) {
		t.Fail()
	}

	q = Quaternion{1, 2, 3, 4}
	r = q.Normalize()
	if !float.Equal(r.Dot(r), 1) {
		t.Fail()
	}
}

func TestAxisAngle(t *testing.T) {
	q := QuatAxisAngle(1, 0, 0, float64(math.Pi)/2)

	t.Log("Rot X,pi/2:", q)

	v := Quaternion{0, 1, 0, 1}
	v = q.Mul(v).Mul(q.Conjugate())

	t.Log("Vector:", v)

	if !float.Equal(v.X, 0) {
		t.Error("X")
	}
	if !float.Equal(v.Y, 0) {
		t.Error("Y")
	}
	if !float.Equal(v.Z, 1) {
		t.Error("Z")
	}
	if !float.Equal(v.W, 1) {
		t.Error("W")
	}

	vx, vy, vz := q.Rotate(0, 1, 0)
	if !float.Equal(vx, 0) {
		t.Error("x")
	}
	if !float.Equal(vy, 0) {
		t.Error("y")
	}
	if !float.Equal(vz, 1) {
		t.Error("z")
	}
}
