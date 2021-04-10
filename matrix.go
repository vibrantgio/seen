package seen

import (
	"math"

	"github.com/reactivego/seen/float"
	"github.com/reactivego/seen/quat"
)

// Matrix is a 4x4 homogenous matrix use to transform the points of an object
// from 'Object Space' to 'Parent Space'. When separated into constituent
// scale, rotate and translate components is best constructed in the
// order T * R * S where transformations are applied right to left.
// So first scale then rotate and finally translate. This order is arbitrary
// but it feels the most natural. An object in 'Object Space' is first scaled,
// then it is rotated around the 'Object Space' origin and finally it is
// translated in the 'Parent Space' to its final location.
type Matrix [4][4]float64

// IdentityMatrix is a transform that does a 1:1 scale, no rotation and
// no translation.
var IdentityMatrix = Matrix{
	{1, 0, 0, 0},
	{0, 1, 0, 0},
	{0, 0, 1, 0},
	{0, 0, 0, 1},
}

// Matrixer is an interface implemented by objects that can return
// a homogeneous 4x4 matrix (M[row][col]) as an [4][4]float array.
type Matrixer interface {
	Matrix() [4][4]float64
}

// M converts a Matrixer to a Matrix.
func M(m Matrixer) Matrix { return Matrix(m.Matrix()) }

// Frustum will return a matrix capabable of projecting points inside the cube specified
// by l,r,b,t,n,f = -r,r,-t,t,-n,-f to clip coordinates. For all valid (non clipped)
// coordinates the following condition holds: -wc < xc,yc,zc < wc
// All coordinates for which this condition doesn't hold need to be clipped.
// To get from clip space coordinates to native device coordinates divide the xc,yc,zc by wc.
// So xn,yn,zn = xc/wc,yc/wc,zc/wc will give coordinates in de range [-1,1]. These need to be
// mapped via a viewport to screen coordinates.
func Frustum(r, t, n, f float64) Matrix {
	return Matrix{
		{n / r, 0, 0, 0},
		{0, n / t, 0, 0},
		{0, 0, (f + n) / (n - f), 2 * f * n / (n - f)},
		{0, 0, -1, 0},
	}
}

func Ortho(r, t, n, f float64) Matrix {
	return Matrix{
		{1 / r, 0, 0, 0},
		{0, 1 / t, 0, 0},
		{0, 0, 2 / (n - f), (f + n) / (n - f)},
		{0, 0, 0, 1},
	}
}

func Perspective(fovy, aspect, near, far float64) Matrix {
	t := math.Tan(0.5*fovy) * near
	r := t * aspect
	return Frustum(r, t, near, far)
}

func Scale(sx, sy, sz float64) Matrix {
	return Matrix{
		{sx, 0, 0, 0},
		{0, sy, 0, 0},
		{0, 0, sz, 0},
		{0, 0, 0, 1},
	}
}
func Rotate(q quat.Quaternion) Matrix {
	return M(q)
}

func Translate(tx, ty, tz float64) Matrix {
	return Matrix{
		{1, 0, 0, tx},
		{0, 1, 0, ty},
		{0, 0, 1, tz},
		{0, 0, 0, 1},
	}
}

func (m Matrix) Scale(sx, sy, sz float64) Matrix {
	return m.Mul(Matrix{
		{sx, 0, 0, 0},
		{0, sy, 0, 0},
		{0, 0, sz, 0},
		{0, 0, 0, 1}})
}

func (m Matrix) Rotate(q quat.Quaternion) Matrix {
	return m.Mul(M(q))
}

func (m Matrix) Translate(tx, ty, tz float64) Matrix {
	return m.Mul(Matrix{
		{1, 0, 0, tx},
		{0, 1, 0, ty},
		{0, 0, 1, tz},
		{0, 0, 0, 1}})
}

// Multiply 4x4 matrices.
func (l Matrix) Mul(r Matrix) Matrix {
	return Matrix{
		{l[0][0]*r[0][0] + l[0][1]*r[1][0] + l[0][2]*r[2][0] + l[0][3]*r[3][0], l[0][0]*r[0][1] + l[0][1]*r[1][1] + l[0][2]*r[2][1] + l[0][3]*r[3][1], l[0][0]*r[0][2] + l[0][1]*r[1][2] + l[0][2]*r[2][2] + l[0][3]*r[3][2], l[0][0]*r[0][3] + l[0][1]*r[1][3] + l[0][2]*r[2][3] + l[0][3]*r[3][3]},
		{l[1][0]*r[0][0] + l[1][1]*r[1][0] + l[1][2]*r[2][0] + l[1][3]*r[3][0], l[1][0]*r[0][1] + l[1][1]*r[1][1] + l[1][2]*r[2][1] + l[1][3]*r[3][1], l[1][0]*r[0][2] + l[1][1]*r[1][2] + l[1][2]*r[2][2] + l[1][3]*r[3][2], l[1][0]*r[0][3] + l[1][1]*r[1][3] + l[1][2]*r[2][3] + l[1][3]*r[3][3]},
		{l[2][0]*r[0][0] + l[2][1]*r[1][0] + l[2][2]*r[2][0] + l[2][3]*r[3][0], l[2][0]*r[0][1] + l[2][1]*r[1][1] + l[2][2]*r[2][1] + l[2][3]*r[3][1], l[2][0]*r[0][2] + l[2][1]*r[1][2] + l[2][2]*r[2][2] + l[2][3]*r[3][2], l[2][0]*r[0][3] + l[2][1]*r[1][3] + l[2][2]*r[2][3] + l[2][3]*r[3][3]},
		{l[3][0]*r[0][0] + l[3][1]*r[1][0] + l[3][2]*r[2][0] + l[3][3]*r[3][0], l[3][0]*r[0][1] + l[3][1]*r[1][1] + l[3][2]*r[2][1] + l[3][3]*r[3][1], l[3][0]*r[0][2] + l[3][1]*r[1][2] + l[3][2]*r[2][2] + l[3][3]*r[3][2], l[3][0]*r[0][3] + l[3][1]*r[1][3] + l[3][2]*r[2][3] + l[3][3]*r[3][3]},
	}
}

func (m Matrix) Transform(x, y, z, w float64) (rx, ry, rz, rw float64) {
	rx = m[0][0]*x + m[0][1]*y + m[0][2]*z + m[0][3]*w
	ry = m[1][0]*x + m[1][1]*y + m[1][2]*z + m[1][3]*w
	rz = m[2][0]*x + m[2][1]*y + m[2][2]*z + m[2][3]*w
	rw = m[3][0]*x + m[3][1]*y + m[3][2]*z + m[3][3]*w
	return
}

func (l Matrix) Equal(r Matrix) bool {
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if !float.Equal(l[i][j], r[i][j]) {
				return false
			}
		}
	}
	return true
}

func (m Matrix) Transpose() (t Matrix) {
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			t[j][i] = m[i][j]
		}
	}
	return
}

func (m Matrix) Determinant() float64 {
	return m[0][0]*m[1][1]*m[2][2]*m[3][3] -
		m[0][0]*m[1][1]*m[2][3]*m[3][2] -
		m[0][0]*m[1][2]*m[2][1]*m[3][3] +
		m[0][0]*m[1][2]*m[2][3]*m[3][1] +
		m[0][0]*m[1][3]*m[2][1]*m[3][2] -
		m[0][0]*m[1][3]*m[2][2]*m[3][1] -
		m[0][1]*m[1][0]*m[2][2]*m[3][3] +
		m[0][1]*m[1][0]*m[2][3]*m[3][2] +
		m[0][1]*m[1][2]*m[2][0]*m[3][3] -
		m[0][1]*m[1][2]*m[2][3]*m[3][0] -
		m[0][1]*m[1][3]*m[2][0]*m[3][2] +
		m[0][1]*m[1][3]*m[2][2]*m[3][0] +
		m[0][2]*m[1][0]*m[2][1]*m[3][3] -
		m[0][2]*m[1][0]*m[2][3]*m[3][1] -
		m[0][2]*m[1][1]*m[2][0]*m[3][3] +
		m[0][2]*m[1][1]*m[2][3]*m[3][0] +
		m[0][2]*m[1][3]*m[2][0]*m[3][1] -
		m[0][2]*m[1][3]*m[2][1]*m[3][0] -
		m[0][3]*m[1][0]*m[2][1]*m[3][2] +
		m[0][3]*m[1][0]*m[2][2]*m[3][1] +
		m[0][3]*m[1][1]*m[2][0]*m[3][2] -
		m[0][3]*m[1][1]*m[2][2]*m[3][0] -
		m[0][3]*m[1][2]*m[2][0]*m[3][1] +
		m[0][3]*m[1][2]*m[2][1]*m[3][0]

}
