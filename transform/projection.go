package transform

import "math"

type Projection struct {
	R, T, N, F float64
}

func Frustrum(fovy, aspect, near, far float64) Projection {
	t := math.Tan(0.5*fovy) * near
	r := t * aspect
	return Projection{r, t, near, far}
}

// Perspective will project points inside the cube specified by l,r,b,t,n,f = -r,r,-t,t,-n,-f
// to clip coordinates. For all valid (non clipped) coordinates the following condition
//  holds: -wc < xc,yc,zc < wc
// All coordinates for which this condition doesn't hold need to be clipped.
// To get from clip space coordinates to native device coordinates divide the xc,yc,zc by wc.
// So xn,yn,zn = xc/wc,yc/wc,zc/wc will give coordinates in de range [-1,1]. These need to be
// mapped via a viewport to screen coordinates.
func (p Projection) Perspective(xe, ye, ze float64) (xc, yc, zc, wc float64) {
	r, t, n, f := p.R, p.T, p.N, p.F
	xc = xe * n / r
	yc = ye * n / t
	zc = (ze*(f+n) + 2*f*n) / (n - f)
	wc = -ze
	return
}

func (p Projection) PerspectiveMat4x4() Mat4x4 {
	r, t, n, f := p.R, p.T, p.N, p.F
	return Mat4x4{
		n / r, 0, 0, 0,
		0, n / t, 0, 0,
		0, 0, (f + n) / (n - f), 2 * f * n / (n - f),
		0, 0, -1, 0,
	}
}

func (p Projection) Orthographic(xe, ye, ze float64) (xc, yc, zc float64) {
	r, t, n, f := p.R, p.T, p.N, p.F
	xc = xe / r
	yc = ye / t
	zc = (ze*-2 - (f + n)) / (f - n)
	return
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

type AsymmetricProjection struct {
	L, R, B, T, N, F float64
}

func (p AsymmetricProjection) Perspective(xe, ye, ze float64) (xc, yc, zc, wc float64) {
	l, r, b, t, n, f := p.L, p.R, p.B, p.T, p.N, p.F
	xc = (xe*2*n + ze*(r+l)) / (r - l)
	yc = (ye*2*n + ze*(t+b)) / (t - b)
	zc = (ze*(f+n) + 2*f*n) / (n - f)
	wc = -ze
	return
}

func (p AsymmetricProjection) Orthographic(xe, ye, ze float64) (xc, yc, zc float64) {
	l, r, b, t, n, f := p.L, p.R, p.B, p.T, p.N, p.F
	xc = (xe*2 - (r + l)) / (r - l)
	yc = (ye*2 - (t + b)) / (t - b)
	zc = (ze*-2 - (f + n)) / (f - n)
	return
}
