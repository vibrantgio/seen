package dualquat

import (
	"fmt"

	"github.com/reactivego/seen/float"
	"github.com/reactivego/seen/quat"
)

// DualQuaternion is a concatenated Translation*Rotation.
// A point transformed with it, first has rotation applied
// and then translation. This behavior is identical to the
// order a homogenous matrix would perform the operation.
type DualQuaternion struct{ Real, Dual quat.Quaternion }

var Identity = DualQuaternion{quat.Identity, quat.Zero}

// TransRot is shorthand for Translate(x,y,z).Rotate(r)
func TransRot(x, y, z float64, r quat.Quaternion) DualQuaternion {
	return DualQuaternion{r, quat.Q(x, y, z, 0).Mul(r).Scale(0.5)}
}

func Translate(x, y, z float64) DualQuaternion {
	return DualQuaternion{quat.Identity, quat.Q(x, y, z, 0).Scale(0.5)}
}

func Rotate(r quat.Quaternion) DualQuaternion {
	return DualQuaternion{r, quat.Zero}
}

func (dq DualQuaternion) Translate(x, y, z float64) DualQuaternion {
	return dq.Mul(Translate(x, y, z))
}

func (dq DualQuaternion) Rotate(q quat.Quaternion) DualQuaternion {
	return dq.Mul(Rotate(q))
}

func (dq DualQuaternion) Rotation() quat.Quaternion {
	return dq.Real
}

// Translation will extract the translation vector from the Dual quaternion.
// To extract the translation 19 muls, 12 adds are needed
func (dq DualQuaternion) Translation() (x, y, z float64) {
	t := dq.Dual.Mul(dq.Real.Conjugate()).Scale(2.0)
	return t.X, t.Y, t.Z
}

func (dq DualQuaternion) Conjugate() DualQuaternion {
	return DualQuaternion{dq.Real.Conjugate(), dq.Dual.Conjugate()}
}

func (lhs DualQuaternion) Add(rhs DualQuaternion) DualQuaternion {
	return DualQuaternion{lhs.Real.Add(rhs.Real), lhs.Dual.Add(rhs.Dual)}
}

func (dq DualQuaternion) Scale(scale float64) DualQuaternion {
	return DualQuaternion{dq.Real.Scale(scale), dq.Dual.Scale(scale)}
}

func (lhs DualQuaternion) Dot(rhs DualQuaternion) float64 {
	return lhs.Real.Dot(rhs.Real)
}

// Mul multiplies two dual quaternions. This takes 3 quaternion multiplications and 1 quaternion addition.
// In total this takes 48 muls and 40 adds
func (lhs DualQuaternion) Mul(rhs DualQuaternion) DualQuaternion {
	//	lhs.Real*rhs.Real, lhs.Real*rhs.Dual + lhs.Dual*rhs.Real
	return DualQuaternion{lhs.Real.Mul(rhs.Real), lhs.Real.Mul(rhs.Dual).Add(lhs.Dual.Mul(rhs.Real))}
}

func (dq DualQuaternion) Normalize() DualQuaternion {
	magnitude := dq.Real.Length()
	//detect near zero magnitude
	if float.Equal(magnitude, 0) {
		return Identity
	}
	return dq.Scale(1.0 / magnitude)
}

// Transform will transform a vector in the space described by the dual quaternion and transform
// that into the parent space.
// This takes 37 muls, 27 adds
// By comparison a homogenous matrix transform takes 9 muls and 9 adds but needs to have the Matrix
// extracted first.
func (dq DualQuaternion) Transform(x, y, z float64) (float64, float64, float64) {
	vx, vy, vz := dq.Real.Transform(x, y, z) // 18 muls and 12 adds
	dx, dy, dz := dq.Translation()           // 19 muls, 12 adds
	return vx + dx, vy + dy, vz + dz         // 3 adds
}

// Matrix will return a matrix with 4 rows and 4 columns, the top left 3x3 matrix
// contains the rotation and the top right 3x1 vector contains the translation.
// It takes 38 muls, 28 adds to derive a homogenous matrix from the dual quaternion.
func (dq DualQuaternion) Matrix() [16]float64 {
	// Returns the homogeneous 3D rotation matrix corresponding to the Real quaternion.
	x, y, z, w := dq.Real.X, dq.Real.Y, dq.Real.Z, dq.Real.W
	// Pre-multiply resused products
	xx, yy, zz := x*x, y*y, z*z // 3 muls
	xy, wz := x*y, w*z          // 2 muls
	xz, wy := x*z, w*y          // 2 muls
	yz, wx := y*z, w*x          // 2 muls
	// Returns the translation corresponding to the Dual quaternion
	tx, ty, tz := dq.Translation() // 20 muls, 12 adds
	// Return a homogenous matrix
	return [16]float64{
		1 - 2*(yy+zz), 2 * (xy - wz), 2 * (xz + wy), tx, // 3 muls, 4 adds
		2 * (xy + wz), 1 - 2*(xx+zz), 2 * (yz - wx), ty, // 3 muls, 4 adds
		2 * (xz - wy), 2 * (yz + wx), 1 - 2*(xx+yy), tz, // 3 muls, 4 adds
		0, 0, 0, 1,
	}
}

func (lhs DualQuaternion) Equal(rhs DualQuaternion) bool {
	return lhs.Real.Equal(rhs.Real) && lhs.Dual.Equal(rhs.Dual)
}

func (dq DualQuaternion) String() string {
	return fmt.Sprintf("{%v, %v}", dq.Real, dq.Dual)
}
