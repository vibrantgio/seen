package matrix

import (
	"math"

	"github.com/vibrantgio/seen/float"
)

// Matrix is a 4x4 homogeneous matrix used to transform an object's points
// from 'Object Space' to 'Parent Space'. When broken into scale, rotation,
// and translation components, it is best constructed in the order T * R * S,
// where transformations are applied from right to left. The object is scaled
// first, then rotated around its origin in 'Object Space', and finally
// translated in 'Parent Space' to its final position. This sequence is
// largely a matter of choice, but it is often considered the most intuitive.
type Matrix [4][4]float64

// Identity is a matrix that does a 1:1 scale, no rotation and no translation.
var Identity = Matrix{
	{1, 0, 0, 0},
	{0, 1, 0, 0},
	{0, 0, 1, 0},
	{0, 0, 0, 1},
}

func Translate(tx, ty, tz float64) Matrix {
	return Matrix{
		{1, 0, 0, tx},
		{0, 1, 0, ty},
		{0, 0, 1, tz},
		{0, 0, 0, 1},
	}
}

// RotateX applies a rotation about the X axis. `theta` is measured in Radians.
func RotateX(theta float64) Matrix {
	st, ct := math.Sincos(theta)
	return Matrix{
		{1, 0, 0, 0},
		{0, ct, -st, 0},
		{0, st, ct, 0},
		{0, 0, 0, 1},
	}
}

// RotateY applies a rotation about the Y axis. `theta` is measured in Radians.
func RotateY(theta float64) Matrix {
	st, ct := math.Sincos(theta)
	return Matrix{
		{ct, 0, st, 0},
		{0, 1, 0, 0},
		{-st, 0, ct, 0},
		{0, 0, 0, 1},
	}
}

// RotateZ applies a rotation about the Z axis. `theta` is measured in Radians.
func RotateZ(theta float64) Matrix {
	st, ct := math.Sincos(theta)
	return Matrix{
		{ct, -st, 0, 0},
		{st, ct, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
}

func Scale(sx, sy, sz float64) Matrix {
	return Matrix{
		{sx, 0, 0, 0},
		{0, sy, 0, 0},
		{0, 0, sz, 0},
		{0, 0, 0, 1},
	}
}

func (m Matrix) Translate(tx, ty, tz float64) Matrix {
	return m.Mul(Matrix{
		{1, 0, 0, tx},
		{0, 1, 0, ty},
		{0, 0, 1, tz},
		{0, 0, 0, 1},
	})
}

// RotateX applies a rotation about the X axis. `theta` is measured in Radians.
func (m Matrix) RotateX(theta float64) Matrix {
	st, ct := math.Sincos(theta)
	return m.Mul(Matrix{
		{1, 0, 0, 0},
		{0, ct, -st, 0},
		{0, st, ct, 0},
		{0, 0, 0, 1},
	})
}

// RotateY applies a rotation about the Y axis. `theta` is measured in Radians.
func (m Matrix) RotateY(theta float64) Matrix {
	st, ct := math.Sincos(theta)
	return m.Mul(Matrix{
		{ct, 0, st, 0},
		{0, 1, 0, 0},
		{-st, 0, ct, 0},
		{0, 0, 0, 1},
	})
}

// RotateZ applies a rotation about the Z axis. `theta` is measured in Radians.
func (m Matrix) RotateZ(theta float64) Matrix {
	st, ct := math.Sincos(theta)
	return m.Mul(Matrix{
		{ct, -st, 0, 0},
		{st, ct, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	})
}

func (m Matrix) Scale(sx, sy, sz float64) Matrix {
	return m.Mul(Matrix{
		{sx, 0, 0, 0},
		{0, sy, 0, 0},
		{0, 0, sz, 0},
		{0, 0, 0, 1},
	})
}

// Multiply 4x4 matrices.
// Uses 64 multiplications and 48 additions.
func (l Matrix) Mul(r Matrix) Matrix {
	return Matrix{
		{l[0][0]*r[0][0] + l[0][1]*r[1][0] + l[0][2]*r[2][0] + l[0][3]*r[3][0], l[0][0]*r[0][1] + l[0][1]*r[1][1] + l[0][2]*r[2][1] + l[0][3]*r[3][1], l[0][0]*r[0][2] + l[0][1]*r[1][2] + l[0][2]*r[2][2] + l[0][3]*r[3][2], l[0][0]*r[0][3] + l[0][1]*r[1][3] + l[0][2]*r[2][3] + l[0][3]*r[3][3]},
		{l[1][0]*r[0][0] + l[1][1]*r[1][0] + l[1][2]*r[2][0] + l[1][3]*r[3][0], l[1][0]*r[0][1] + l[1][1]*r[1][1] + l[1][2]*r[2][1] + l[1][3]*r[3][1], l[1][0]*r[0][2] + l[1][1]*r[1][2] + l[1][2]*r[2][2] + l[1][3]*r[3][2], l[1][0]*r[0][3] + l[1][1]*r[1][3] + l[1][2]*r[2][3] + l[1][3]*r[3][3]},
		{l[2][0]*r[0][0] + l[2][1]*r[1][0] + l[2][2]*r[2][0] + l[2][3]*r[3][0], l[2][0]*r[0][1] + l[2][1]*r[1][1] + l[2][2]*r[2][1] + l[2][3]*r[3][1], l[2][0]*r[0][2] + l[2][1]*r[1][2] + l[2][2]*r[2][2] + l[2][3]*r[3][2], l[2][0]*r[0][3] + l[2][1]*r[1][3] + l[2][2]*r[2][3] + l[2][3]*r[3][3]},
		{l[3][0]*r[0][0] + l[3][1]*r[1][0] + l[3][2]*r[2][0] + l[3][3]*r[3][0], l[3][0]*r[0][1] + l[3][1]*r[1][1] + l[3][2]*r[2][1] + l[3][3]*r[3][1], l[3][0]*r[0][2] + l[3][1]*r[1][2] + l[3][2]*r[2][2] + l[3][3]*r[3][2], l[3][0]*r[0][3] + l[3][1]*r[1][3] + l[3][2]*r[2][3] + l[3][3]*r[3][3]},
	}
}

// Transform a vector.
// Uses 16 multiplications and 12 additions.
func (m Matrix) Transform(vx, vy, vz, vw float64) (x, y, z, w float64) {
	x = m[0][0]*vx + m[0][1]*vy + m[0][2]*vz + m[0][3]*vw
	y = m[1][0]*vx + m[1][1]*vy + m[1][2]*vz + m[1][3]*vw
	z = m[2][0]*vx + m[2][1]*vy + m[2][2]*vz + m[2][3]*vw
	w = m[3][0]*vx + m[3][1]*vy + m[3][2]*vz + m[3][3]*vw
	return
}

// Transform a vector.
// Uses 9 multiplications and 9 additions.
func (m Matrix) Transform3(vx, vy, vz float64) (x, y, z float64) {
	x = m[0][0]*vx + m[0][1]*vy + m[0][2]*vz + m[0][3]
	y = m[1][0]*vx + m[1][1]*vy + m[1][2]*vz + m[1][3]
	z = m[2][0]*vx + m[2][1]*vy + m[2][2]*vz + m[2][3]
	return
}

func (l Matrix) Equal(r Matrix) bool {
	for i := range 4 {
		for j := range 4 {
			if !float.Equal(l[i][j], r[i][j]) {
				return false
			}
		}
	}
	return true
}

func (m Matrix) Transpose() (t Matrix) {
	for i := range 4 {
		for j := range 4 {
			t[j][i] = m[i][j]
		}
	}
	return
}

func (m Matrix) Determinant() float64 {
	return (m[0][0]*m[1][1]*(m[2][2]*m[3][3]-m[3][2]*m[2][3]) -
		m[0][0]*m[1][2]*(m[2][1]*m[3][3]-m[3][1]*m[2][3]) +
		m[0][0]*m[1][3]*(m[2][1]*m[3][2]-m[3][1]*m[2][2])) -
		(m[0][1]*m[1][0]*(m[2][2]*m[3][3]-m[3][2]*m[2][3]) -
			m[0][1]*m[1][2]*(m[2][0]*m[3][3]-m[3][0]*m[2][3]) +
			m[0][1]*m[1][3]*(m[2][0]*m[3][2]-m[3][0]*m[2][2])) +
		(m[0][2]*m[1][0]*(m[2][1]*m[3][3]-m[3][1]*m[2][3]) -
			m[0][2]*m[1][1]*(m[2][0]*m[3][3]-m[3][0]*m[2][3]) +
			m[0][2]*m[1][3]*(m[2][0]*m[3][1]-m[3][0]*m[2][1])) -
		(m[0][3]*m[1][0]*(m[2][1]*m[3][2]-m[3][1]*m[2][2]) -
			m[0][3]*m[1][1]*(m[2][0]*m[3][2]-m[3][0]*m[2][2]) +
			m[0][3]*m[1][2]*(m[2][0]*m[3][1]-m[3][0]*m[2][1]))
}

func (m Matrix) Invert() (inv Matrix, ok bool) {
	minor := func(r0, r1, r2, c0, c1, c2 int) float64 {
		return m[r0][c0]*(m[r1][c1]*m[r2][c2]-m[r2][c1]*m[r1][c2]) -
			m[r0][c1]*(m[r1][c0]*m[r2][c2]-m[r2][c0]*m[r1][c2]) +
			m[r0][c2]*(m[r1][c0]*m[r2][c1]-m[r2][c0]*m[r1][c1])
	}

	det := m[0][0]*minor(1, 2, 3, 1, 2, 3) -
		m[0][1]*minor(1, 2, 3, 0, 2, 3) +
		m[0][2]*minor(1, 2, 3, 0, 1, 3) -
		m[0][3]*minor(1, 2, 3, 0, 1, 2)

	// Singularity must be judged relative to the matrix's own scale, not by
	// an absolute epsilon: a uniform scale s contributes s^4 to det, so a
	// perfectly invertible view matrix with a small prescale (e.g. 1/2200
	// per axis) has det ~ 1e-11 and a raw float.Equal(det, 0) misreads it
	// as singular. Hadamard's inequality bounds |det| by the product of the
	// row norms; a well-conditioned matrix stays within a modest factor of
	// that bound, while a rank-deficient one collapses to ~machine epsilon
	// of it. Zero row norms (and NaN/Inf) are singular outright.
	scale := 1.0
	for r := range 4 {
		n := math.Sqrt(m[r][0]*m[r][0] + m[r][1]*m[r][1] + m[r][2]*m[r][2] + m[r][3]*m[r][3])
		scale *= n
	}
	if ok = scale > 0 && !math.IsNaN(det) && !math.IsInf(det, 0) && math.Abs(det) > 1e-12*scale; !ok {
		return
	}

	inv[0][0] = minor(1, 2, 3, 1, 2, 3) / det
	inv[0][1] = -minor(0, 2, 3, 1, 2, 3) / det
	inv[0][2] = minor(0, 1, 3, 1, 2, 3) / det
	inv[0][3] = -minor(0, 1, 2, 1, 2, 3) / det
	inv[1][0] = -minor(1, 2, 3, 0, 2, 3) / det
	inv[1][1] = minor(0, 2, 3, 0, 2, 3) / det
	inv[1][2] = -minor(0, 1, 3, 0, 2, 3) / det
	inv[1][3] = minor(0, 1, 2, 0, 2, 3) / det
	inv[2][0] = minor(1, 2, 3, 0, 1, 3) / det
	inv[2][1] = -minor(0, 2, 3, 0, 1, 3) / det
	inv[2][2] = minor(0, 1, 3, 0, 1, 3) / det
	inv[2][3] = -minor(0, 1, 2, 0, 1, 3) / det
	inv[3][0] = -minor(1, 2, 3, 0, 1, 2) / det
	inv[3][1] = minor(0, 2, 3, 0, 1, 2) / det
	inv[3][2] = -minor(0, 1, 3, 0, 1, 2) / det
	inv[3][3] = minor(0, 1, 2, 0, 1, 2) / det
	return
}
