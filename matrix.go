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
type Matrix [16]float64

// IdentityMatrix is a transform that does a 1:1 scale, no rotation and
// no translation.
var IdentityMatrix = Matrix{
	1, 0, 0, 0,
	0, 1, 0, 0,
	0, 0, 1, 0,
	0, 0, 0, 1,
}

// Matrixer is an interface implemented by objects that can return
// a homogeneous 4x4 matrix (per row) as an [16]float array.
type Matrixer interface {
	Matrix() [16]float64
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
		n / r, 0, 0, 0,
		0, n / t, 0, 0,
		0, 0, (f + n) / (n - f), 2 * f * n / (n - f),
		0, 0, -1, 0,
	}
}

func Ortho(r, t, n, f float64) Matrix {
	return Matrix{
		1 / r, 0, 0, 0,
		0, 1 / t, 0, 0,
		0, 0, 2 / (n - f), (f + n) / (n - f),
		0, 0, 0, 1,
	}
}

func Perspective(fovy, aspect, near, far float64) Matrix {
	t := math.Tan(0.5*fovy) * near
	r := t * aspect
	return Frustum(r, t, near, far)
}

func Scale(sx, sy, sz float64) Matrix {
	return Matrix{
		sx, 0, 0, 0,
		0, sy, 0, 0,
		0, 0, sz, 0,
		0, 0, 0, 1,
	}
}
func Rotate(q quat.Quaternion) Matrix {
	return M(q)
}

func Translate(tx, ty, tz float64) Matrix {
	return Matrix{
		1, 0, 0, tx,
		0, 1, 0, ty,
		0, 0, 1, tz,
		0, 0, 0, 1,
	}
}

func (m Matrix) Scale(sx, sy, sz float64) Matrix {
	return m.Mul(Matrix{
		sx, 0, 0, 0,
		0, sy, 0, 0,
		0, 0, sz, 0,
		0, 0, 0, 1})
}

func (m Matrix) Rotate(q quat.Quaternion) Matrix {
	return m.Mul(M(q))
}

func (m Matrix) Translate(tx, ty, tz float64) Matrix {
	return m.Mul(Matrix{
		1, 0, 0, tx,
		0, 1, 0, ty,
		0, 0, 1, tz,
		0, 0, 0, 1})
}

// Multiply 4x4 matrices.
func (l Matrix) Mul(r Matrix) Matrix {
	return Matrix{
		l[0]*r[0] + l[1]*r[4] + l[2]*r[8] + l[3]*r[12], l[0]*r[1] + l[1]*r[5] + l[2]*r[9] + l[3]*r[13], l[0]*r[2] + l[1]*r[6] + l[2]*r[10] + l[3]*r[14], l[0]*r[3] + l[1]*r[7] + l[2]*r[11] + l[3]*r[15],
		l[4]*r[0] + l[5]*r[4] + l[6]*r[8] + l[7]*r[12], l[4]*r[1] + l[5]*r[5] + l[6]*r[9] + l[7]*r[13], l[4]*r[2] + l[5]*r[6] + l[6]*r[10] + l[7]*r[14], l[4]*r[3] + l[5]*r[7] + l[6]*r[11] + l[7]*r[15],
		l[8]*r[0] + l[9]*r[4] + l[10]*r[8] + l[11]*r[12], l[8]*r[1] + l[9]*r[5] + l[10]*r[9] + l[11]*r[13], l[8]*r[2] + l[9]*r[6] + l[10]*r[10] + l[11]*r[14], l[8]*r[3] + l[9]*r[7] + l[10]*r[11] + l[11]*r[15],
		l[12]*r[0] + l[13]*r[4] + l[14]*r[8] + l[15]*r[12], l[12]*r[1] + l[13]*r[5] + l[14]*r[9] + l[15]*r[13], l[12]*r[2] + l[13]*r[6] + l[14]*r[10] + l[15]*r[14], l[12]*r[3] + l[13]*r[7] + l[14]*r[11] + l[15]*r[15],
	}
}

func (m Matrix) Transform(x, y, z, w float64) (rx, ry, rz, rw float64) {
	rx = m[0]*x + m[1]*y + m[2]*z + m[3]*w
	ry = m[4]*x + m[5]*y + m[6]*z + m[7]*w
	rz = m[8]*x + m[9]*y + m[10]*z + m[11]*w
	rw = m[12]*x + m[13]*y + m[14]*z + m[15]*w
	return
}

func (l Matrix) Equal(r Matrix) bool {
	for i, li := range l {
		if !float.Equal(li, r[i]) {
			return false
		}
	}
	return true
}

func (m Matrix) Transpose() (t Matrix) {
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			t[j*4+i] = m[i*4+j]
		}
	}
	return
}
