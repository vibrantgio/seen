package quaternion

import (
	"fmt"

	"github.com/vibrantgio/seen/float"
	"github.com/vibrantgio/seen/matrix"
)

// Dual is a concatenated Translation*Rotation.
// A point transformed with it, first has rotation applied
// and then translation. This behavior is identical to the
// order a homogenous matrix would perform the operation.
type Dual struct{ Real, Dual Quat }

// DualIdentity returns the Dual Quaternion that induces zero translation and zero
// rotation when multiplied with another Dual Quaternion.
var DualIdentity = Dual{Real: Identity}

// DQ is shorthand for Translate(x,y,z).Rotate(r).
func DQ(tx, ty, tz float64, r Quat) Dual {
	return Dual{r, Q(tx, ty, tz, 0).Mul(r).Scale(0.5)}
}

func (dq Dual) Translate(tx, ty, tz float64) Dual {
	return dq.Mul(Dual{Identity, Q(tx, ty, tz, 0).Scale(0.5)})
}

func (dq Dual) Rotate(r Quat) Dual {
	return dq.Mul(Dual{Real: r})
}

func (dq Dual) Rotation() Quat {
	return dq.Real
}

// Translation will extract the translation vector from the Dual quaternion.
// To extract the translation 19 muls, 12 adds are needed
func (dq Dual) Translation() (x, y, z float64) {
	t := dq.Dual.Mul(dq.Real.Conjugate()).Scale(2.0)
	return t.X, t.Y, t.Z
}

func (dq Dual) Conjugate() Dual {
	return Dual{dq.Real.Conjugate(), dq.Dual.Conjugate()}
}

func (lhs Dual) Add(rhs Dual) Dual {
	return Dual{lhs.Real.Add(rhs.Real), lhs.Dual.Add(rhs.Dual)}
}

func (dq Dual) Scale(scale float64) Dual {
	return Dual{dq.Real.Scale(scale), dq.Dual.Scale(scale)}
}

func (lhs Dual) Dot(rhs Dual) float64 {
	return lhs.Real.Dot(rhs.Real)
}

// Mul multiplies two dual quaternions.
// Uses 3 quaternion muls and 1 quaternion add.
// In total this uses 48 muls and 40 adds.
func (lhs Dual) Mul(rhs Dual) Dual {
	//	lhs.Real*rhs.Real, lhs.Real*rhs.Dual + lhs.Dual*rhs.Real
	return Dual{lhs.Real.Mul(rhs.Real), lhs.Real.Mul(rhs.Dual).Add(lhs.Dual.Mul(rhs.Real))}
}

func (dq Dual) Normalize() Dual {
	magnitude := dq.Real.Length()
	//detect near zero magnitude
	if float.Equal(magnitude, 0) {
		return Dual{Real: Identity}
	}
	return dq.Scale(1.0 / magnitude)
}

// Transform will transform a vector in the space described by the dual
// quaternion and transform that into the parent space. This takes 37 muls, 27
// adds. By comparison a homogenous matrix transform takes 9 muls and 9 adds but
// needs to have the Matrix extracted first.
func (dq Dual) Transform(x, y, z float64) (float64, float64, float64) {
	vx, vy, vz := dq.Real.Transform(x, y, z) // 18 muls and 12 adds
	dx, dy, dz := dq.Translation()           // 19 muls, 12 adds
	return vx + dx, vy + dy, vz + dz         // 3 adds
}

// Matrix will return a matrix with 4 rows and 4 columns, the top left 3x3
// matrix contains the rotation and the top right 3x1 vector contains the
// translation. It takes 38 muls, 28 adds to derive a homogenous matrix from the
// dual quaternion.
func (dq Dual) Matrix() matrix.Matrix {
	// Returns the homogeneous 3D rotation matrix corresponding to the Real quaternion.
	x, y, z, w := dq.Real.X, dq.Real.Y, dq.Real.Z, dq.Real.W
	// Pre-multiply reused products
	xx, yy, zz := x*x, y*y, z*z // 3 muls
	xy, wz := x*y, w*z          // 2 muls
	xz, wy := x*z, w*y          // 2 muls
	yz, wx := y*z, w*x          // 2 muls
	// Returns the translation corresponding to the Dual quaternion
	tx, ty, tz := dq.Translation() // 20 muls, 12 adds
	// Return a homogenous matrix
	return [4][4]float64{
		{1 - 2*(yy+zz), 2 * (xy - wz), 2 * (xz + wy), tx}, // 3 muls, 4 adds
		{2 * (xy + wz), 1 - 2*(xx+zz), 2 * (yz - wx), ty}, // 3 muls, 4 adds
		{2 * (xz - wy), 2 * (yz + wx), 1 - 2*(xx+yy), tz}, // 3 muls, 4 adds
		{0, 0, 0, 1},
	}
}

func (lhs Dual) Equal(rhs Dual) bool {
	return lhs.Real.Equal(rhs.Real) && lhs.Dual.Equal(rhs.Dual)
}

func (dq Dual) String() string {
	return fmt.Sprintf("{%v, %v}", dq.Real, dq.Dual)
}
