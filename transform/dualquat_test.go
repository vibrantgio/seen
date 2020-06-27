package transform

import (
	"math"
	"testing"
	"github.com/reactivego/seen/float"
)

/*
DualQuaternion_c dq0 = new DualQuaternion_c(
	Quaternion.CreateFromYawPitchRoll(1,2,3),
	new Vector3(10,30,90) );

DualQuaternion_c dq1 = new DualQuaternion_c(
	Quaternion.CreateFromYawPitchRoll(-1,3,2),
	new Vector3(30,40, 190 ) );

DualQuaternion_c dq2 = new DualQuaternion_c(
	Quaternion.CreateFromYawPitchRoll(2,3,1.5f),
	new Vector3(5,20, 66 ) );

DualQuaternion_c dq = dq0 * dq1 * dq2;

Matrix dqToMatrix = DualQuaternion_c.DualQuaternionToMatrix( dq );
Matrix m0 = Matrix.CreateFromYawPitchRoll(1,2,3) * Matrix.CreateTranslation( 10, 30, 90 );
Matrix m1 = Matrix.CreateFromYawPitchRoll(-1,3,2) * Matrix.CreateTranslation( 30, 40, 190 );
Matrix m2 = Matrix.CreateFromYawPitchRoll(2,3,1.5f) * Matrix.CreateTranslation( 5, 20, 66 );
Matrix m = m0 * m1 * m2;

*/

func TestCreationDQ(t *testing.T) {

	d := IdentDualQuaternion

	if d == nil {
		t.Fail()
	}
}

func TestInitializationDQ(t *testing.T) {

	r := (&Quaternion{1, 2, 3, 4}).Normalize()
	dq := MakeDualQuatRXYZ(r, 5, 6, 7)
	x, y, z := dq.Translation()
	if !float.Equal(x, 5) {
		t.Fail()
	}
	if !float.Equal(y, 6) {
		t.Fail()
	}
	if !float.Equal(z, 7) {
		t.Fail()
	}
}

func TestTransformRotateDQ(t *testing.T) {

	dq := MakeDualQuatRXYZ(MakeQuatAxisAngle(1, 0, 0, math.Pi/2), 0, 0, 0)
	x, y, z := dq.Transform(0, 1, 0)

	t.Log(dq, x, y, z)

	if !float.Equal(x, 0) {
		t.Fail()
	}
	if !float.Equal(y, 0) {
		t.Fail()
	}
	if !float.Equal(z, 1) {
		t.Fail()
	}
}

func TestTransformTranslateDQ(t *testing.T) {

	dq := MakeDualQuatRXYZ(IdentQuaternion, 1, 2, 3)
	x, y, z := dq.Transform(4, 5, 6)

	if !float.Equal(x, 5) {
		t.Fail()
	}
	if !float.Equal(y, 7) {
		t.Fail()
	}
	if !float.Equal(z, 9) {
		t.Fail()
	}
}

func TestTransformCombinedDQ(t *testing.T) {

	dq := MakeDualQuatRXYZ(MakeQuatAxisAngle(1, 0, 0, math.Pi/2), 1, 2, 3)
	x, y, z := dq.Transform(4, 5, 6)

	// Point 4,5,6 in object space has been rotated around x axis by 90 degrees and then translated with vector 1,2,3

	t.Log(dq, x, y, z)

	if !float.Equal(x, 5) {
		t.Fail()
	}
	if !float.Equal(y, -4) {
		t.Fail()
	}
	if !float.Equal(z, 8) {
		t.Fail()
	}
}

func TestTransformStacked(t *testing.T) {

	// dq0 transforms from world space to view space
	dq0 := MakeDualQuatRXYZ(MakeQuatAxisAngle(0, 0, 1, math.Pi), 100, 100, 100)
	// dq1 transforms from intermediate space to world space
	dq1 := MakeDualQuatRXYZ(MakeQuatAxisAngle(1, 0, 0, -math.Pi/2), 20, 20, 20)
	// dq2 transforms from object space to intermediate space
	dq2 := MakeDualQuatRXYZ(MakeQuatAxisAngle(1, 0, 0, math.Pi/2), 1, 2, 3)

	// dq12 transforms from object space to world space.
	dq12 := dq1.Mul(dq2)
	// dq02 transforms from object space to view space.
	dq02 := dq0.Mul(dq12)

	// To determine what the  x,y,z values should be, we analyze the transformations taking place.

	// Object to Intermediate space
	// Point 4,5,6 in object space has been rotated around x axis by 90 degrees resulting in 4,-6,5
	// The result is then translated with vector 1,2,3 resulting in 5,-4,8
	x, y, z := dq2.Transform(4, 5, 6)
	if !float.Equal(x, 5) {
		t.Fail()
	}
	if !float.Equal(y, -4) {
		t.Fail()
	}
	if !float.Equal(z, 8) {
		t.Fail()
	}

	// Intermediate to World space
	// The previous result is rotated back 90 degrees around the x axis resulting in 5,8,4
	// The result is then translated with vector 20,20,20 resulting in 25,28,24
	x, y, z = dq1.Transform(x, y, z)
	if !float.Equal(x, 25) {
		t.Fail()
	}
	if !float.Equal(y, 28) {
		t.Fail()
	}
	if !float.Equal(z, 24) {
		t.Fail()
	}
	// Also verify that stacked object space to world space transform gives identical result
	x, y, z = dq12.Transform(4, 5, 6)
	if !float.Equal(x, 25) {
		t.Fail()
	}
	if !float.Equal(y, 28) {
		t.Fail()
	}
	if !float.Equal(z, 24) {
		t.Fail()
	}

	// World to View space
	// The previous result is rotated around z axis 180 degrees resulting in sign change for x and y values -25,-28,24
	// The result is then translated with vector 100,100,100 resulting in 75,72,124
	x, y, z = dq0.Transform(x, y, z)
	if !float.Equal(x, 75) {
		t.Fail()
	}
	if !float.Equal(y, 72) {
		t.Fail()
	}
	if !float.Equal(z, 124) {
		t.Fail()
	}
	// Also verify that stacked object space to view space transform gives identical result
	x, y, z = dq02.Transform(4, 5, 6)
	if !float.Equal(x, 75) {
		t.Fail()
	}
	if !float.Equal(y, 72) {
		t.Fail()
	}
	if !float.Equal(z, 124) {
		t.Fail()
	}
}
