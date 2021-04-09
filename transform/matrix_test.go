package transform

import (
	"math"
	"testing"

	"github.com/reactivego/seen/float"
)

func TestMatrixFromDualQuat(t *testing.T) {

	dq := DualQuatRXYZ(QuatAxisAngle(1, 2, 3, math.Pi/2.0), 4, 5, 6)
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

	dq1 := DualQuatRXYZ(QuatAxisAngle(1, 2, 3, math.Pi/2.0), 4, 5, 6)
	dq2 := DualQuatRXYZ(QuatAxisAngle(0, 1, 0, math.Pi/4.0), 10, 11, 12)
	dq3 := DualQuatRXYZ(QuatAxisAngle(0, 1, 1, math.Pi/3.0), 100, 110, 120)
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
