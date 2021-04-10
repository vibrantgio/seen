package seen

import (
	"math"
	"testing"

	"github.com/reactivego/seen/dualquat"
	"github.com/reactivego/seen/float"
	"github.com/reactivego/seen/quat"
)

func TestMatrixFromDualQuat(t *testing.T) {

	dq := dualquat.TransRot(4, 5, 6, quat.AxisAngle(1, 2, 3, math.Pi/2.0))
	m := M(dq)

	dx, dy, dz := dq.Transform(7, 8, 9)
	mx, my, mz, _ := m.Transform(7, 8, 9, 1)

	if !float.EqualPairs(dx, mx, dy, my, dz, mz) {
		t.Logf("Exp: {%.4v,%.4v,%.4v}", dx, dy, dz)
		t.Logf("Got: {%.4v,%.4v,%.4v}", mx, my, mz)
		t.Fail()
	}
}

func TestMatrixDualQuatOperationOrder(t *testing.T) {

	r := quat.AxisAngle(1, 2, 3, math.Pi/2.0)

	R := dualquat.Rotate(r)
	T := dualquat.Translate(4, 5, 6)

	R_T := R.Mul(T)
	MR_MT := M(R).Mul(M(T))

	if !M(R_T).Equal(MR_MT) {
		t.Error("M(R_T) != MR_MT")
	}

	T_R := T.Mul(R)
	TR := dualquat.TransRot(4, 5, 6, r)

	if !T_R.Equal(TR) {
		t.Error("T_R != TR")
	}

	TRr := T.Rotate(r)

	if !T_R.Equal(TRr) {
		t.Error("T_R != TRr")
	}
}

func TestMatrixMultiplication(t *testing.T) {

	dq1 := dualquat.TransRot(4, 5, 6, quat.AxisAngle(1, 2, 3, math.Pi/2.0))
	dq2 := dualquat.TransRot(10, 11, 12, quat.AxisAngle(0, 1, 0, math.Pi/4.0))
	dq3 := dualquat.TransRot(100, 110, 120, quat.AxisAngle(0, 1, 1, math.Pi/3.0))
	dq := dq1.Mul(dq2).Mul(dq3)

	m1 := M(dq1)
	m2 := M(dq2)
	m3 := M(dq3)
	m := m1.Mul(m2).Mul(m3)

	dx, dy, dz := dq.Transform(7, 8, 9)
	mx, my, mz, mw := m.Transform(7, 8, 9, 1)

	if !float.EqualPairs(dx, mx, dy, my, dz, mz, 1.0, mw) {
		t.Errorf("Exp: {%.4v,%.4v,%.4v,%.4v}\nGot: {%.4v,%.4v,%.4v,%.4v}",
			dx, dy, dz, 1.0, mx, my, mz, mw)
	}
}

func TestFrustum(t *testing.T) {
	var (
		width  float64 = 4000
		height float64 = 3000
		near   float64 = 100
		far    float64 = 200
		right          = 0.5 * width
		left           = -right
		top            = 0.5 * height
		bottom         = -top
	)
	m := Frustum(right, top, near, far)

	// Check that zero value for x,y,z returns the correct values
	x, y, z, w := m.Transform(0, 0, 0, 1)
	if !float.EqualPairs(x, 0, y, 0, z, -400, w, 0) {
		t.Fail()
	}

	// Check that positive edges are equal to w value
	x, y, z, w = m.Transform(right, top, -near, 1)
	if !float.EqualPairs(x, w, y, w, z, -w) {
		t.Fail()
	}

	// check that negative edges are equal to negative w value
	x, y, z, w = m.Transform(left, bottom, -near, 1)
	if !float.EqualPairs(x, -w, y, -w, z, -w) {
		t.Fail()
	}

	// check that clipping works by nudging outside the edge
	x, y, _, w = m.Transform(left-1.0, bottom-1.0, -near, 1)
	if x >= -w {
		t.Fail()
	}
	if y >= -w {
		t.Fail()
	}
}
