package transform

import "math"

// Projection is a 3D to 2D `Matrix` transformation.
// A projection assumes the camera is located at (0,0,0).
type Projection struct {
	R, T, N, F float64
}

func Frustum(r, t, n, f float64) Matrix {
	return Projection{r, t, n, f}.PerspectiveMatrix()
}

func Ortho(r, t, n, f float64) Matrix {
	return Projection{r, t, n, f}.OrthographicMatrix()
}

func Perspective(fovy, aspect, near, far float64) Matrix {
	t := math.Tan(0.5*fovy) * near
	r := t * aspect
	return Projection{r, t, near, far}.PerspectiveMatrix()
}

// PerspectiveMatrix will return a matrix capabable of projecting points inside the cube
// specified by l,r,b,t,n,f = -r,r,-t,t,-n,-f to clip coordinates. For all valid (non clipped)
// coordinates the following condition holds: -wc < xc,yc,zc < wc
// All coordinates for which this condition doesn't hold need to be clipped.
// To get from clip space coordinates to native device coordinates divide the xc,yc,zc by wc.
// So xn,yn,zn = xc/wc,yc/wc,zc/wc will give coordinates in de range [-1,1]. These need to be
// mapped via a viewport to screen coordinates.
func (p Projection) PerspectiveMatrix() Matrix {
	r, t, n, f := p.R, p.T, p.N, p.F
	return Matrix{
		n / r, 0, 0, 0,
		0, n / t, 0, 0,
		0, 0, (f + n) / (n - f), 2 * f * n / (n - f),
		0, 0, -1, 0,
	}
}

func (p Projection) OrthographicMatrix() Matrix {
	r, t, n, f := p.R, p.T, p.N, p.F
	return Matrix{
		1 / r, 0, 0, 0,
		0, 1 / t, 0, 0,
		0, 0, 2 / (n - f), (f + n) / (n - f),
		0, 0, 0, 1,
	}
}
