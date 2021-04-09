package transform

import (
	"math"
	"testing"

	"github.com/reactivego/seen/float"
	"github.com/reactivego/seen/quat"
)

func TestMatrixFromDualQuat(t *testing.T) {

	dq := DualQuatRXYZ(quat.AxisAngle(1, 2, 3, math.Pi/2.0), 4, 5, 6)
	m := dq.Matrix()

	dx, dy, dz := dq.Transform(7, 8, 9)
	mx, my, mz, _ := m.Transform(7, 8, 9, 1)

	t.Log("dx", dx, "dy", dy, "dz", dz)
	t.Log("mx", mx, "my", my, "mz", mz)

	if !float.EqualPairs(mx, dx, my, dy, mz, dz) {
		t.Fail()
	}
}

func TestMatrixMultiplication(t *testing.T) {

	dq1 := DualQuatRXYZ(quat.AxisAngle(1, 2, 3, math.Pi/2.0), 4, 5, 6)
	dq2 := DualQuatRXYZ(quat.AxisAngle(0, 1, 0, math.Pi/4.0), 10, 11, 12)
	dq3 := DualQuatRXYZ(quat.AxisAngle(0, 1, 1, math.Pi/3.0), 100, 110, 120)
	dq := dq1.Mul(dq2).Mul(dq3)

	m1 := dq1.Matrix()
	m2 := dq2.Matrix()
	m3 := dq3.Matrix()
	m := m1.Mul(m2).Mul(m3)

	dx, dy, dz := dq.Transform(7, 8, 9)
	mx, my, mz, mw := m.Transform(7, 8, 9, 1)

	t.Log("dx", dx, "dy", dy, "dz", dz)
	t.Log("mx", mx, "my", my, "mz", mz, "mw", mw)

	if !float.EqualPairs(mx, dx, my, dy, mz, dz, mw, 1.0) {
		t.Fail()
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

func TestPerspective(t *testing.T) {
	fovy := 135.0 * (math.Pi / 180.0)
	p := Perspective(fovy, 4.0/3.0, 100.0, 200.0)
	t.Log(p)

	// TODO: Actually test something
}
