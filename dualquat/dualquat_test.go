package dualquat

import (
	"math"
	"testing"

	"github.com/reactivego/seen/float"
	"github.com/reactivego/seen/quat"
)

func TestInitializationDQ(t *testing.T) {

	r := quat.Q(1, 2, 3, 4).Normalize()
	TR := TransRot(5, 6, 7, r)
	x, y, z := TR.Translation()

	if !float.EqualPairs(x, 5, y, 6, z, 7) {
		t.Errorf("Exp: {5,6,7}\nGot: {%v,%v,%v}", x, y, z)
	}

	T_R := Translate(5, 6, 7).Rotate(r)
	x, y, z = T_R.Translation()

	if !float.EqualPairs(x, 5, y, 6, z, 7) {
		t.Errorf("Exp: {5,6,7}\nGot: {%v,%v,%v}", x, y, z)
	}

	if !TR.Equal(T_R) {
		t.Error("TR != T_R")
	}
}

func TestTransformRotateDQ(t *testing.T) {

	TR := TransRot(0, 0, 0, quat.AxisAngle(1, 0, 0, math.Pi/2))
	x, y, z := TR.Transform(0, 1, 0)

	if !float.EqualPairs(x, 0, y, 0, z, 1) {
		t.Errorf("Exp: {0,0,1}\nGot: {%v,%v,%v}", x, y, z)
	}
}

func TestTransformTranslateDQ(t *testing.T) {

	RT := TransRot(1, 2, 3, quat.Identity)
	x, y, z := RT.Transform(4, 5, 6)

	if !float.EqualPairs(x, 5, y, 7, z, 9) {
		t.Logf("expected {5,6,7}, got {%v,%v,%v}", x, y, z)
		t.Fail()
	}
}

func TestTransformCombinedDQ(t *testing.T) {

	// Rotate point {4,5,6} in object space around x axis by 90 degrees,
	// then translate by vector {1,2,3}
	TR := TransRot(1, 2, 3, quat.AxisAngle(1, 0, 0, math.Pi/2))
	x, y, z := TR.Transform(4, 5, 6)

	if !float.EqualPairs(x, 5, y, -4, z, 8) {
		t.Errorf("Exp: {5,-4,8}\nGot: {%v,%v,%v}", x, y, z)
	}
}

func TestTransformStacked(t *testing.T) {

	// dq0 transforms from world space to view space
	dq0 := TransRot(100, 100, 100, quat.AxisAngle(0, 0, 1, math.Pi))
	// dq1 transforms from intermediate space to world space
	dq1 := TransRot(20, 20, 20, quat.AxisAngle(1, 0, 0, -math.Pi/2))
	// dq2 transforms from object space to intermediate space
	dq2 := TransRot(1, 2, 3, quat.AxisAngle(1, 0, 0, math.Pi/2))

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
