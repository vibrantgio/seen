package transform

import "math"

// Projection is a 3D to 2D `Matrix` transformation.
// A projection assumes the camera is located at (0,0,0).
type Projection struct {
	R, T, N, F float64
}

func Frustum(r, t, n, f float64) Mat4x4 {
	return Projection{r, t, n, f}.PerspectiveMat4x4()
}

func Ortho(r, t, n, f float64) Mat4x4 {
	return Projection{r, t, n, f}.OrthographicMat4x4()
}

func Perspective(fovy, aspect, near, far float64) Mat4x4 {
	t := math.Tan(0.5*fovy) * near
	r := t * aspect
	return Projection{r, t, near, far}.PerspectiveMat4x4()
}

// PerspectiveMat4x4 will return a matrix capabable of projecting points inside the cube
// specified by l,r,b,t,n,f = -r,r,-t,t,-n,-f to clip coordinates. For all valid (non clipped)
// coordinates the following condition holds: -wc < xc,yc,zc < wc
// All coordinates for which this condition doesn't hold need to be clipped.
// To get from clip space coordinates to native device coordinates divide the xc,yc,zc by wc.
// So xn,yn,zn = xc/wc,yc/wc,zc/wc will give coordinates in de range [-1,1]. These need to be
// mapped via a viewport to screen coordinates.
func (p Projection) PerspectiveMat4x4() Mat4x4 {
	r, t, n, f := p.R, p.T, p.N, p.F
	return Mat4x4{
		n / r, 0, 0, 0,
		0, n / t, 0, 0,
		0, 0, (f + n) / (n - f), 2 * f * n / (n - f),
		0, 0, -1, 0,
	}
}

func (p Projection) OrthographicMat4x4() Mat4x4 {
	r, t, n, f := p.R, p.T, p.N, p.F
	return Mat4x4{
		1 / r, 0, 0, 0,
		0, 1 / t, 0, 0,
		0, 0, 2 / (n - f), (f + n) / (n - f),
		0, 0, 0, 1,
	}
}
