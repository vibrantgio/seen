package seen

import (
	"math"
	"testing"

	"github.com/reactivego/seen/dualquat"
	"github.com/reactivego/seen/float"
	"github.com/reactivego/seen/quat"
)

func TestMatrixRotationOrthonormal(t *testing.T) {

	R := M(quat.RotX(0.1 * math.Pi).RotY(0.1 * math.Pi).RotZ(0.1 * math.Pi))
	RT := R.Transpose()

	I := R.Mul(RT)

	if !I.Equal(IdentityMatrix) {
		t.Log("Exp: IdentityMatrix")
		t.Errorf("Got: %.4v", I)
	}
}

func TestMatrixInvert(t *testing.T) {

	R := M(quat.RotX(0.1*math.Pi).RotY(0.1*math.Pi).RotZ(0.1*math.Pi)).Scale(1, 2, 3)
	if Ri, ok := R.Invert(); !ok {
		t.Error("matrix not invertable")
	} else {

		I := R.Mul(Ri)

		if !I.Equal(IdentityMatrix) {
			t.Log("R * Ri")
			t.Log("  Exp: IdentityMatrix")
			t.Errorf("  Got: %.4v", I)
		}

		det := R.Determinant()
		if float.Equal(det, 0.0) {
			t.Errorf("Determinant != 0 but got %.4v", det)
		}
	}
}

func TestMatrixTRS(t *testing.T) {

	QR := quat.RotX(0.1 * math.Pi).RotY(0.1 * math.Pi).RotZ(0.1 * math.Pi)

	QTRS := M(dualquat.TransRot(1, 2, 3, QR)).Scale(2.0, 2.0, 2.0)
	MTRS := Translate(1, 2, 3).Rotate(QR).Scale(2.0, 2.0, 2.0)

	if !QTRS.Equal(MTRS) {
		t.Logf("QTRS: %.5v", QTRS)
		t.Logf("MTRS: %.5v", MTRS)
		t.Error("did not expect QTRS != MTRS")
	}
}

func TestMatrixExtractScale(t *testing.T) {
	// Scale can be extracted from a homogeneous matrix under
	// the following assumptions:
	// 1. R has to be orthonormal => Rt * R = R * Rt = I  (Rt is R transposed)
	// 2. S components have to be >=0

	// Given a full blown homegeneous transformation matrix TRS
	// Note that TRS => T * RS, we can extract RS by zeroing out
	// TRS03,TRS13 and TRS23 giving us RS.

	// Now note that (RS)t * RS => St * (Rt * R) * S = St * I * S = St * S
	// Also note that St == S as the scaling is on the diagonal which
	// is invariant under transpose. So St * S = SS the values at SS00,
	// SS11 and SS22 give the squared scale components SX*SX,SY*SY,SZ*SZ.
	// Use square roots to get the original SX, SY and SZ.

	R := M(quat.RotX(0.1 * math.Pi).RotY(0.1 * math.Pi).RotZ(0.1 * math.Pi))
	S := Scale(2, 3, 4)

	RS := R.Mul(S)
	StRt := RS.Transpose()
	StS := StRt.Mul(RS)

	sx, sy, sz := math.Sqrt(StS[0][0]), math.Sqrt(StS[1][1]), math.Sqrt(StS[2][2])

	if !float.EqualPairs(sx, 2, sy, 3, sz, 4) {
		t.Errorf("EqualPairs:\nExp: {2,3,4}\nGot: {%.5v,%.5v,%.5v}", sx, sy, sz)
	}
}

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
