package quat

import (
	"fmt"
	"math"

	"github.com/reactivego/seen/float"
)

type Quaternion struct {
	X, Y, Z, W float64
}

var Identity = Quaternion{0, 0, 0, 1}

var Zero = Quaternion{}

func Q(x, y, z, w float64) Quaternion {
	return Quaternion{x, y, z, w}
}

func AxisAngle(x, y, z, angle float64) Quaternion {
	// determine length of axis angle so we can normalize.
	m := math.Sqrt(x*x + y*y + z*z)
	// filter out degenerate axis.
	if float.Equal(m, 0) {
		return Identity
	}
	s, c := math.Sincos(angle / 2)
	return Quaternion{s * x / m, s * y / m, s * z / m, c}
}

func RotX(angle float64) Quaternion {
	s, c := math.Sincos(angle / 2)
	return Quaternion{s, 0.0, 0.0, c}
}

func RotY(angle float64) Quaternion {
	s, c := math.Sincos(angle / 2)
	return Quaternion{0.0, s, 0.0, c}
}

func RotZ(angle float64) Quaternion {
	s, c := math.Sincos(angle / 2)
	return Quaternion{0.0, 0.0, s, c}
}

func (q Quaternion) String() string {
	return fmt.Sprintf("{X=%v,Y=%v,Z=%v,W=%v}", q.X, q.Y, q.Z, q.W)
}

func (lhs Quaternion) Equal(rhs Quaternion) bool {
	return float.Equal(lhs.X, rhs.X) && float.Equal(lhs.Y, rhs.Y) && float.Equal(lhs.Z, rhs.Z) && float.Equal(lhs.W, rhs.W)
}

func (q Quaternion) Conjugate() Quaternion {
	return Quaternion{-q.X, -q.Y, -q.Z, q.W}
}

// Add returns the sum of two quaternions. This takes 4 adds.
func (lhs Quaternion) Add(rhs Quaternion) Quaternion {
	return Quaternion{lhs.X + rhs.X, lhs.Y + rhs.Y, lhs.Z + rhs.Z, lhs.W + rhs.W}
}

func (q Quaternion) Scale(scale float64) Quaternion {
	return Quaternion{q.X * scale, q.Y * scale, q.Z * scale, q.W * scale}
}

// Dot returns the quaternion dot product (inner product) of the target (q) and r.
// (For two normalized quaternions, this will be 1 if they’re equal, -1 if they’re opposite and 0 if they’re perpendicular.)
func (lhs Quaternion) Dot(rhs Quaternion) float64 {
	return lhs.X*rhs.X + lhs.Y*rhs.Y + lhs.Z*rhs.Z + lhs.W*rhs.W
}

func (q Quaternion) Length() float64 {
	return math.Sqrt(q.Dot(q))
}

func (q Quaternion) Normalize() Quaternion {
	magnitude := q.Length()
	//detect near zero magnitude
	if float.Equal(magnitude, 0) {
		return Identity
	}
	return q.Scale(1 / magnitude)
}

// Mul calculates the Hamilton product of two quaternions. This can be seen as a rotation.
// Note that Multiplication is NOT commutative, meaning q1.Mul(q2) does not necessarily
// equal q2.Mul(q1).
// This operation takes 16 muls and 12 adds or an alternative implemnentation can
// do it in 9 muls and 27 adds. It's not known whether adds on modern x86 cpu's are still
// faster than muls.
func (lhs Quaternion) Mul(rhs Quaternion) Quaternion {
	// 16 muls and 12 adds
	return Quaternion{
		lhs.W*rhs.X + lhs.X*rhs.W + lhs.Y*rhs.Z - lhs.Z*rhs.Y,
		lhs.W*rhs.Y - lhs.X*rhs.Z + lhs.Y*rhs.W + lhs.Z*rhs.X,
		lhs.W*rhs.Z + lhs.X*rhs.Y - lhs.Y*rhs.X + lhs.Z*rhs.W,
		lhs.W*rhs.W - lhs.X*rhs.X - lhs.Y*rhs.Y - lhs.Z*rhs.Z,
	}

	// Alternative implementation
	// 9 muls, 27 adds
	/*	ww := (lhs.Z + lhs.X) * (rhs.X + rhs.Y)
		yy := (lhs.W - lhs.Y) * (rhs.W + rhs.Z)
		zz := (lhs.W + lhs.Y) * (rhs.W - rhs.Z)
		xx := ww + yy + zz
		qq := 0.5 * (xx + (lhs.Z-lhs.X)*(rhs.X-rhs.Y))

		x := qq - xx + (lhs.X+lhs.W)*(rhs.X+rhs.W)
		y := qq - yy + (lhs.W-lhs.X)*(rhs.Y+rhs.Z)
		z := qq - zz + (lhs.Z+lhs.Y)*(rhs.W-rhs.X)
		w := qq - ww + (lhs.Z-lhs.Y)*(rhs.Y-rhs.Z)
		return &Quaternion{x, y, z, w}
	*/
}

// RotX multiplies a quaternion with a Rotation around the x-axis. q' = qqX
func (q Quaternion) RotX(angle float64) Quaternion {
	s, c := math.Sincos(angle / 2)
	return q.Mul(Quaternion{s, 0.0, 0.0, c})
}

// RotY multiplies a quaternion with a Rotation around the y-axis. q' = qqY
func (q Quaternion) RotY(angle float64) Quaternion {
	s, c := math.Sincos(angle / 2)
	return q.Mul(Quaternion{0.0, s, 0.0, c})
}

// RotZ multiplies a quaternion with a Rotation around the z-axis. q' = qqZ
func (q Quaternion) RotZ(angle float64) Quaternion {
	s, c := math.Sincos(angle / 2)
	return q.Mul(Quaternion{0.0, 0.0, s, c})
}

// Rotate will perform q*v*q' on the passed in vector.
// This takes 18 muls, 12 adds to compute.
func (q Quaternion) Transform(x, y, z float64) (rx, ry, rz float64) {
	cross := func(ax, ay, az, bx, by, bz float64) (rx, ry, rz float64) {
		// cross product is implemented with 6 muls, 3 adds
		return ay*bz - az*by, az*bx - ax*bz, ax*by - ay*bx
	}

	// Standard implementation
	// 32 muls, 24 adds
	/*	t := q.Mul(&Quaternion{x, y, z, 0}).Mul(q.Conjugate())
		return t.X, t.Y, t.Z
	*/

	// Alternative implementation
	// 18 muls and 12 adds
	tx, ty, tz := cross(2*q.X, 2*q.Y, 2*q.Z, x, y, z) // 9 muls, 3 adds
	ux, uy, uz := cross(q.X, q.Y, q.Z, tx, ty, tz)    // 6 muls, 3 adds
	rx = x + q.W*tx + ux                              // 1 mul, 2 adds
	ry = y + q.W*ty + uy                              // 1 mul, 2 adds
	rz = z + q.W*tz + uz                              // 1 mul, 2 adds
	return
}

// Matrix will return a matrix with 4 rows and 4 columns, the top left 3x3 matrix
// contains the rotation. Computing the 4x4 homogeneous matrix from the quaternion
// takes 18 muls and 12 adds
func (q Quaternion) Matrix() [16]float64 {
	// Returns the homogeneous 3D rotation matrix corresponding to the quaternion.
	x, y, z, w := q.X, q.Y, q.Z, q.W
	// Pre-multiply resused products
	xx, yy, zz := x*x, y*y, z*z // 3 muls
	xy, wz := x*y, w*z          // 2 muls
	xz, wy := x*z, w*y          // 2 muls
	yz, wx := y*z, w*x          // 2 muls
	// Return a homogenous matrix
	return [16]float64{
		1 - 2*(yy+zz), 2 * (xy - wz), 2 * (xz + wy), 0, // 3 muls, 4 adds
		2 * (xy + wz), 1 - 2*(xx+zz), 2 * (yz - wx), 0, // 3 muls, 4 adds
		2 * (xz - wy), 2 * (yz + wx), 1 - 2*(xx+yy), 0, // 3 muls, 4 adds
		0, 0, 0, 1,
	}
}
