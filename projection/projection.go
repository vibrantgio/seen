package projection

import (
	"math"

	"github.com/vibrantgio/seen/matrix"
)

// Frustum will return a matrix capabable of projecting points inside the cube specified
// by l,r,b,t,n,f = -r,r,-t,t,-n,-f to clip coordinates. For all valid (non clipped)
// coordinates the following condition holds: -wc < xc,yc,zc < wc
// All coordinates for which this condition doesn't hold need to be clipped.
// To get from clip space coordinates to native device coordinates divide the xc,yc,zc by wc.
// So xn,yn,zn = xc/wc,yc/wc,zc/wc will give coordinates in de range [-1,1]. These need to be
// mapped via a viewport to screen coordinates.
func Frustum(r, t, n, f float64) matrix.Matrix {
	return matrix.Matrix{
		{n / r, 0, 0, 0},
		{0, n / t, 0, 0},
		{0, 0, (f + n) / (n - f), 2 * f * n / (n - f)},
		{0, 0, -1, 0},
	}
}

func Perspective(fovy, aspect, near, far float64) matrix.Matrix {
	t := math.Tan(0.5*fovy) * near
	r := t * aspect
	return Frustum(r, t, near, far)
}

// DefaultPerspective is the default projection matrix using Frustum.
var DefaultPerspective = Frustum(1, 1, 1, 100)

func Ortho(r, t, n, f float64) matrix.Matrix {
	return matrix.Matrix{
		{1 / r, 0, 0, 0},
		{0, 1 / t, 0, 0},
		{0, 0, 2 / (n - f), (f + n) / (n - f)},
		{0, 0, 0, 1},
	}
}

// DefaultOrthographic is the default projection matrix using Ortho.
var DefaultOrthographic = Ortho(1, 1, 1, 100)
