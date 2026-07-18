package matrix_test

import (
	"math"
	"testing"

	"github.com/vibrantgio/seen/float"
	"github.com/vibrantgio/seen/matrix"
)

func TestMatrixRotationOrthonormal(t *testing.T) {
	MR := matrix.RotateX(0.1 * math.Pi).RotateY(0.1 * math.Pi).RotateZ(0.1 * math.Pi)
	RT := MR.Transpose()
	I := MR.Mul(RT)
	if !I.Equal(matrix.Identity) {
		t.Log("Exp: IdentityMatrix")
		t.Errorf("Got: %.4v", I)
	}
}

func TestMatrixInvert(t *testing.T) {
	MR := matrix.RotateX(0.1 * math.Pi).RotateY(0.1 * math.Pi).RotateZ(0.1 * math.Pi)
	MRS := MR.Scale(1, 2, 3)
	if Ri, ok := MRS.Invert(); !ok {
		t.Error("matrix not invertable")
	} else {

		I := MRS.Mul(Ri)

		if !I.Equal(matrix.Identity) {
			t.Log("R * Ri")
			t.Log("  Exp: IdentityMatrix")
			t.Errorf("  Got: %.4v", I)
		}

		det := MRS.Determinant()
		if float.Equal(det, 0.0) {
			t.Errorf("Determinant != 0 but got %.4v", det)
		}
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
	MR := matrix.RotateX(0.1 * math.Pi).RotateY(0.1 * math.Pi).RotateZ(0.1 * math.Pi)

	MS := matrix.Scale(2, 3, 4)

	MRxMS := MR.Mul(MS)
	MRxMS_T := MRxMS.Transpose()
	StS := MRxMS_T.Mul(MRxMS)

	sx, sy, sz := math.Sqrt(StS[0][0]), math.Sqrt(StS[1][1]), math.Sqrt(StS[2][2])

	if !float.EqualPairs(sx, 2, sy, 3, sz, 4) {
		t.Errorf("EqualPairs:\nExp: {2,3,4}\nGot: {%.5v,%.5v,%.5v}", sx, sy, sz)
	}
}

// TestMatrixInvertSmallScale: a view matrix with a small uniform scale
// (1/2200 per axis, as Scene.FitCenter sets Camera.Norm for a fixed reference
// distance) has det ~ 1e-11 — far below any absolute epsilon, yet perfectly
// invertible. The singularity test must be relative to the matrix's own scale.
func TestMatrixInvertSmallScale(t *testing.T) {
	s := 1.0 / 2200.0
	view := matrix.Scale(s, s, s).Translate(0, 0, -2200)
	inv, ok := view.Invert()
	if !ok {
		t.Fatal("small-prescale view matrix reported singular")
	}
	if I := view.Mul(inv); !I.Equal(matrix.Identity) {
		t.Errorf("view * view^-1 != I, got %.4v", I)
	}
	// The eye recovery this exists for: the preimage of the origin.
	x, y, z := inv.Transform3(0, 0, 0)
	if math.Abs(x) > 1e-9 || math.Abs(y) > 1e-9 || math.Abs(z-2200) > 1e-6 {
		t.Errorf("eye = (%v, %v, %v), want (0, 0, 2200)", x, y, z)
	}
}

// TestMatrixInvertSingular: genuinely rank-deficient matrices must still be
// rejected under the scale-relative test.
func TestMatrixInvertSingular(t *testing.T) {
	if _, ok := matrix.Scale(0, 1, 1).Invert(); ok {
		t.Error("zero-scale matrix reported invertible")
	}
	collapsed := matrix.Identity
	collapsed[2] = collapsed[1] // two identical rows: det = 0 at unit scale
	if _, ok := collapsed.Invert(); ok {
		t.Error("rank-deficient matrix reported invertible")
	}
}
