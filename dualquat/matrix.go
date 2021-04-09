package transform

import "math"

type Matrix [16]float64

var IdentityMatrix = Matrix{
	1, 0, 0, 0,
	0, 1, 0, 0,
	0, 0, 1, 0,
	0, 0, 0, 1,
}

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
