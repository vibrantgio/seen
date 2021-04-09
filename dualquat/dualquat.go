package transform

import "github.com/reactivego/seen/quat"

type DualQuaternion struct{ Real, Dual quat.Quaternion }

var IdentDualQuaternion = DualQuaternion{quat.Identity, quat.Zero}

func DualQuatXYZ(x, y, z float64) DualQuaternion {
	return DualQuaternion{quat.Identity, quat.Quaternion{0.5 * x, 0.5 * y, 0.5 * z, 0}}
}

func DualQuatR(r quat.Quaternion) DualQuaternion {
	return DualQuaternion{r, quat.Zero}
}

func DualQuatRXYZ(r quat.Quaternion, x, y, z float64) DualQuaternion {
	t := quat.Quaternion{x, y, z, 0}
	d := t.Mul(r).Scale(0.5)
	return DualQuaternion{r, d}
}

func (q DualQuaternion) Rotation() quat.Quaternion {
	return q.Real
}

// Translation will extract the translation vector from the Dual quaternion.
// To extract the translation 19 muls, 12 adds are needed
func (q DualQuaternion) Translation() (x, y, z float64) {
	t := q.Dual.Mul(q.Real.Conjugate())
	return 2.0 * t.X, 2.0 * t.Y, 2.0 * t.Z
}

func (q DualQuaternion) Conjugate() DualQuaternion {
	return DualQuaternion{q.Real.Conjugate(), q.Dual.Conjugate()}
}

func (lhs DualQuaternion) Add(rhs DualQuaternion) DualQuaternion {
	return DualQuaternion{lhs.Real.Add(rhs.Real), lhs.Dual.Add(rhs.Dual)}
}

func (q DualQuaternion) Scale(scale float64) DualQuaternion {
	return DualQuaternion{q.Real.Scale(scale), q.Dual.Scale(scale)}
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

func (q DualQuaternion) Normalize() DualQuaternion {
	// I don't think this Normalize function is actually correct.

	mag := q.Real.Dot(q.Real)
	if mag < 0.000001 {
		return IdentDualQuaternion
	}
	return q.Scale(1.0 / mag)
}

// Transform will transform a vector in the space described by the dual quaternion and transform
// that into the parent space.
// This takes 37 muls, 27 adds
// By comparison a homogenous matrix transform takes 9 muls and 9 adds but needs to have the Matrix
// extracted first.
func (q DualQuaternion) Transform(x, y, z float64) (float64, float64, float64) {
	vx, vy, vz := q.Real.Rotate(x, y, z) // 18 muls and 12 adds
	dx, dy, dz := q.Translation()        // 19 muls, 12 adds
	return vx + dx, vy + dy, vz + dz     // 3 adds
}

// Matrix will return a matrix with 4 rows and 4 columns, the top left 3x3 matrix
// contains the rotation and the top right 3x1 vector contains the translation.
// It takes 38 muls, 28 adds to derive a homogenous matrix from the dual quaternion.
func (q DualQuaternion) Matrix() Matrix {
	// Returns the homogeneous 3D rotation matrix corresponding to the Real quaternion.
	x, y, z, w := q.Real.X, q.Real.Y, q.Real.Z, q.Real.W
	// Pre-multiply resused products
	xx, yy, zz := x*x, y*y, z*z // 3 muls
	xy, wz := x*y, w*z          // 2 muls
	xz, wy := x*z, w*y          // 2 muls
	yz, wx := y*z, w*x          // 2 muls
	// Returns the translation corresponding to the Dual quaternion
	tx, ty, tz := q.Translation() // 20 muls, 12 adds
	// Return a homogenous matrix
	return Matrix{
		1 - 2*(yy+zz), 2 * (xy - wz), 2 * (xz + wy), tx, // 3 muls, 4 adds
		2 * (xy + wz), 1 - 2*(xx+zz), 2 * (yz - wx), ty, // 3 muls, 4 adds
		2 * (xz - wy), 2 * (yz + wx), 1 - 2*(xx+yy), tz, // 3 muls, 4 adds
		0, 0, 0, 1,
	}
}
