package affine

import (
	"testing"

	"github.com/vibrantgio/seen/float"
)

func TestAffineTransformation(t *testing.T) {
	xform := SolveForAffineTransform(ORTHONORMAL_BASIS)
	if !float.EqualPairs(xform.A, 1.0, xform.B, 0.0, xform.C, 0.0, xform.D, -1.0, xform.E, 0.0, xform.F, 0.0) {
		t.Log(ORTHONORMAL_BASIS)
		t.Log(xform)
		t.Fail()
	}
}
