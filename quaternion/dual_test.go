package quaternion_test

import (
	"math"
	"testing"

	"github.com/vibrantgio/seen/float"
	"github.com/vibrantgio/seen/matrix"
	"github.com/vibrantgio/seen/quaternion"
)

func TestInitializationDQ(t *testing.T) {

	r := quaternion.Q(1, 2, 3, 4).Normalize()
	TR := quaternion.DQ(5, 6, 7, r)
	x, y, z := TR.Translation()

	if !float.EqualPairs(x, 5, y, 6, z, 7) {
		t.Errorf("Exp: {5,6,7}\nGot: {%v,%v,%v}", x, y, z)
	}

	T_R := quaternion.DualIdentity.Translate(5, 6, 7).Rotate(r)
	x, y, z = T_R.Translation()

	if !float.EqualPairs(x, 5, y, 6, z, 7) {
		t.Errorf("Exp: {5,6,7}\nGot: {%v,%v,%v}", x, y, z)
	}

	if !TR.Equal(T_R) {
		t.Error("TR != T_R")
	}
}

func TestTransformIdentityDQ(t *testing.T) {
	TR := quaternion.DQ(1, 2, 3, quaternion.AxisAngle(1, 0, 0, math.Pi/2))

	TR1 := quaternion.DualIdentity.Mul(TR)
	if !TR.Equal(TR1) {
		t.Error("TR != TR1")
	}

	TR2 := TR1.Mul(quaternion.DualIdentity)
	if !TR.Equal(TR2) {
		t.Error("TR != TR2")
	}
}

func TestTransformRotateDQ(t *testing.T) {

	TR := quaternion.DQ(0, 0, 0, quaternion.AxisAngle(1, 0, 0, math.Pi/2))
	x, y, z := TR.Transform(0, 1, 0)

	if !float.EqualPairs(x, 0, y, 0, z, 1) {
		t.Errorf("Exp: {0,0,1}\nGot: {%v,%v,%v}", x, y, z)
	}
}

func TestTransformTranslateDQ(t *testing.T) {

	RT := quaternion.DQ(1, 2, 3, quaternion.Identity)
	x, y, z := RT.Transform(4, 5, 6)

	if !float.EqualPairs(x, 5, y, 7, z, 9) {
		t.Logf("expected {5,6,7}, got {%v,%v,%v}", x, y, z)
		t.Fail()
	}
}

func TestTransformCombinedDQ(t *testing.T) {

	// Rotate point {4,5,6} in object space around x axis by 90 degrees,
	// then translate by vector {1,2,3}
	TR := quaternion.DQ(1, 2, 3, quaternion.AxisAngle(1, 0, 0, math.Pi/2))
	x, y, z := TR.Transform(4, 5, 6)

	if !float.EqualPairs(x, 5, y, -4, z, 8) {
		t.Errorf("Exp: {5,-4,8}\nGot: {%v,%v,%v}", x, y, z)
	}
}

func TestTransformStacked(t *testing.T) {

	// dq0 transforms from world space to view space
	dq0 := quaternion.DQ(100, 100, 100, quaternion.AxisAngle(0, 0, 1, math.Pi))
	// dq1 transforms from intermediate space to world space
	dq1 := quaternion.DQ(20, 20, 20, quaternion.AxisAngle(1, 0, 0, -math.Pi/2))
	// dq2 transforms from object space to intermediate space
	dq2 := quaternion.DQ(1, 2, 3, quaternion.AxisAngle(1, 0, 0, math.Pi/2))

	// dq12 transforms from object space to world space.
	dq12 := dq1.Mul(dq2)
	// dq02 transforms from object space to view space.
	dq02 := dq0.Mul(dq12)

	// To determine what the  x,y,z values should be, we analyze the transformations taking place.

	// Object to Intermediate space
	// Point 4,5,6 in object space has been rotated around x axis by 90 degrees resulting in 4,-6,5
	// The result is then translated with vector 1,2,3 resulting in 5,-4,8
	x, y, z := dq2.Transform(4, 5, 6)
	if !float.EqualPairs(x, 5, y, -4, z, 8) {
		t.Fail()
	}

	// Intermediate to World space
	// The previous result is rotated back 90 degrees around the x axis resulting in 5,8,4
	// The result is then translated with vector 20,20,20 resulting in 25,28,24
	x, y, z = dq1.Transform(x, y, z)
	if !float.EqualPairs(x, 25, y, 28, z, 24) {
		t.Fail()
	}
	// Also verify that stacked object space to world space transform gives identical result
	x, y, z = dq12.Transform(4, 5, 6)
	if !float.EqualPairs(x, 25, y, 28, z, 24) {
		t.Fail()
	}

	// World to View space
	// The previous result is rotated around z axis 180 degrees resulting in sign change for x and y values -25,-28,24
	// The result is then translated with vector 100,100,100 resulting in 75,72,124
	x, y, z = dq0.Transform(x, y, z)
	if !float.EqualPairs(x, 75, y, 72, z, 124) {
		t.Fail()
	}
	// Also verify that stacked object space to view space transform gives identical result
	x, y, z = dq02.Transform(4, 5, 6)
	if !float.EqualPairs(x, 75, y, 72, z, 124) {
		t.Fail()
	}
}

func TestMatrixTRS(t *testing.T) {
	QTRS := quaternion.DQ(1, 2, 3, quaternion.RotX(0.1*math.Pi).RotY(0.1*math.Pi).RotZ(0.1*math.Pi)).Matrix().Scale(2.0, 2.0, 2.0)
	MTRS := matrix.Translate(1, 2, 3).RotateX(0.1*math.Pi).RotateY(0.1*math.Pi).RotateZ(0.1*math.Pi).Scale(2.0, 2.0, 2.0)

	if !QTRS.Equal(MTRS) {
		t.Logf("QTRS: %.5v", QTRS)
		t.Logf("MTRS: %.5v", MTRS)
		t.Error("did not expect QTRS != MTRS")
	}
}

func TestMatrixFromDualQuat(t *testing.T) {

	dq := quaternion.DQ(4, 5, 6, quaternion.AxisAngle(1, 2, 3, math.Pi/2.0))
	m := dq.Matrix()

	dx, dy, dz := dq.Transform(7, 8, 9)
	mx, my, mz, _ := m.Transform(7, 8, 9, 1)

	if !float.EqualPairs(dx, mx, dy, my, dz, mz) {
		t.Logf("Exp: {%.4v,%.4v,%.4v}", dx, dy, dz)
		t.Logf("Got: {%.4v,%.4v,%.4v}", mx, my, mz)
		t.Fail()
	}
}

func TestMatrixDualQuatOperationOrder(t *testing.T) {

	r := quaternion.AxisAngle(1, 2, 3, math.Pi/2.0)

	R := quaternion.DQ(0, 0, 0, r)
	T := quaternion.DQ(4, 5, 6, quaternion.Identity)

	M_RT := R.Mul(T).Matrix()
	MR_MT := R.Matrix().Mul(T.Matrix())

	if !M_RT.Equal(MR_MT) {
		t.Error("M(R_T) != MR_MT")
	}

	T_R := T.Mul(R)
	TR := quaternion.DQ(4, 5, 6, r)

	if !T_R.Equal(TR) {
		t.Error("T_R != TR")
	}

	TRr := T.Rotate(r)

	if !T_R.Equal(TRr) {
		t.Error("T_R != TRr")
	}
}

func TestMatrixMultiplication(t *testing.T) {

	dq1 := quaternion.DQ(4, 5, 6, quaternion.AxisAngle(1, 2, 3, math.Pi/2.0))
	dq2 := quaternion.DQ(10, 11, 12, quaternion.AxisAngle(0, 1, 0, math.Pi/4.0))
	dq3 := quaternion.DQ(100, 110, 120, quaternion.AxisAngle(0, 1, 1, math.Pi/3.0))
	dq := dq1.Mul(dq2).Mul(dq3)

	m1 := dq1.Matrix()
	m2 := dq2.Matrix()
	m3 := dq3.Matrix()
	m := m1.Mul(m2).Mul(m3)

	dx, dy, dz := dq.Transform(7, 8, 9)
	mx, my, mz, mw := m.Transform(7, 8, 9, 1)

	if !float.EqualPairs(dx, mx, dy, my, dz, mz, 1.0, mw) {
		t.Errorf("Exp: {%.4v,%.4v,%.4v,%.4v}\nGot: {%.4v,%.4v,%.4v,%.4v}",
			dx, dy, dz, 1.0, mx, my, mz, mw)
	}
}
