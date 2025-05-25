package point

import (
	"math"

	"github.com/vibrantgio/seen/float"
	"github.com/vibrantgio/seen/matrix"
	"github.com/vibrantgio/seen/quaternion"
)

// Point is an object used for specifying the outline of a face.
type Point struct {
	X, Y, Z float64
}

func Pt(x, y, z float64) Point {
	return Point{x, y, z}
}

func P[T int | float64](x, y, z T) Point {
	return Point{float64(x), float64(y), float64(z)}
}

func (v Point) Negated() Point {
	return Point{X: -v.X, Y: -v.Y, Z: -v.Z}
}

func (v Point) Plus(a Point) Point {
	return Point{X: v.X + a.X, Y: v.Y + a.Y, Z: v.Z + a.Z}
}

func (v Point) Minus(a Point) Point {
	return Point{X: v.X - a.X, Y: v.Y - a.Y, Z: v.Z - a.Z}
}

func (v Point) Times(a float64) Point {
	return Point{X: v.X * a, Y: v.Y * a, Z: v.Z * a}
}

func (v Point) DividedBy(a float64) Point {
	return Point{X: v.X / a, Y: v.Y / a, Z: v.Z / a}
}

func (l Point) Dot(r Point) float64 {
	return l.X*r.X + l.Y*r.Y + l.Z*r.Z
}

func (v Point) Lerp(a Point, t float64) Point {
	return v.Plus(a.Minus(v).Times(t))
}

// Length calculates the length of the X,Y,Z component vector
func (p Point) Length() float64 {
	x, y, z := p.X, p.Y, p.Z
	return math.Sqrt(x*x + y*y + z*z)
}

func (v Point) Unit() Point {
	return v.DividedBy(v.Length())
}

func (l Point) Cross(r Point) Point {
	return Point{X: l.Y*r.Z - l.Z*r.Y, Y: l.Z*r.X - l.X*r.Z, Z: l.X*r.Y - l.Y*r.X}
}

// Round rounds the X,Y,Z components to the nearest integer value.
func (p Point) Round() Point {
	return Point{X: math.Floor(p.X + 0.5), Y: math.Floor(p.Y + 0.5), Z: math.Floor(p.Z + 0.5)}
}

// Normalize returns a pointer to a copy of the Point after
// normalizing the copy so the vector length is 1.0
func (p Point) Normalize() Point {
	mag := p.Length()
	return Point{X: p.X / mag, Y: p.Y / mag, Z: p.Z / mag}
}

func (p Point) Perpendicular() Point {
	perp := p.Cross(Point{0, 0, 1})
	if mag := perp.Length(); mag != 0 {
		return perp.DividedBy(mag)
	}
	return p.Cross(Point{1, 0, 0}).Normalize()
}

func (l Point) Equal(r Point) bool {
	return float.EqualPairs(l.X, r.X, l.Y, r.Y, l.Z, r.Z)
}

// Mul returns point p transformed by matrix M, i.e. p' = Mp
func (p Point) Mul(M matrix.Matrix) Point {
	x, y, z, _ := M.Transform(p.X, p.Y, p.Z, 1)
	return Point{X: x, Y: y, Z: z}
}

// PointAngle returns a quaternion representing a rotation around the axis
// defined by the point p. The rotation angle is specified in radians.
func (p Point) PointAngle(angle float64) quaternion.Quat {
	// determine length of axis angle so we can normalize.
	m := math.Sqrt(p.X*p.X + p.Y*p.Y + p.Z*p.Z)
	// filter out degenerate axis.
	if float.Equal(m, 0) {
		return quaternion.Identity
	}
	s, c := math.Sincos(angle / 2)
	return quaternion.Quat{X: s * p.X / m, Y: s * p.Y / m, Z: s * p.Z / m, W: c}

}
