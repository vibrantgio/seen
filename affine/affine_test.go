package affine

import (
	"testing"
	"github.com/reactivego/seen/float"
)

func TestAffineTransformation(t *testing.T) {
	points := ORTHONORMAL_BASIS
	xform := SolveForAffineTransform(points)
	t.Log(points)
	t.Log(xform)
	if !float.EqualPairs(xform.A,1.0, xform.B,0.0, xform.C,0.0, xform.D,-1.0, xform.E,0.0, xform.F,0.0) {
		t.Fail()
	}
}
