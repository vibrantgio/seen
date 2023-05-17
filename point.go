package seen

import (
	"math"

	"github.com/reactivego/seen/float"
)

// Point is an object used for specifying the outline of a Surface.
type Point struct {
	X, Y, Z float64
}

func Pt(x, y, z float64) Point {
	return Point{x, y, z}
}

func P[T int | float64](x, y, z T) Point {
	return Point{float64(x), float64(y), float64(z)}
}

// Normalize returns a pointer to a copy of the Point after
// normalizing the copy so the vector length is 1.0
func (p Point) Normalize() Point {
	mag := p.Length()
	return Point{X: p.X / mag, Y: p.Y / mag, Z: p.Z / mag}
}

func (v Point) Unit() Point {
	return v.DividedBy(v.Length())
}

func (v Point) Negated() Point {
	return Point{X: -v.X, Y: -v.Y, Z: -v.Z}
}

// Subtract returns a pointer to a copy of this Point after
// subtracting the passed in Point r from this copy.
// The W values are not subtracted.
func (l Point) Subtract(r Point) Point {
	return Point{X: l.X - r.X, Y: l.Y - r.Y, Z: l.Z - r.Z}
}

func (v Point) Minus(a Point) Point {
	return Point{X: v.X - a.X, Y: v.Y - a.Y, Z: v.Z - a.Z}
}

// Length calculates the length of the X,Y,Z component vector
func (p Point) Length() float64 {
	x, y, z := p.X, p.Y, p.Z
	return math.Sqrt(x*x + y*y + z*z)
}

func (l Point) Equal(r Point) bool {
	return float.EqualPairs(l.X, r.X, l.Y, r.Y, l.Z, r.Z)
}

// Add will add another Point's X,Y,Z compontents
func (l Point) Add(r Point) Point {
	return Point{X: l.X + r.X, Y: l.Y + r.Y, Z: l.Z + r.Z}
}

func (v Point) Plus(a Point) Point {
	return Point{X: v.X + a.X, Y: v.Y + a.Y, Z: v.Z + a.Z}
}

// Round rounds the X,Y,Z components to the nearest integer value.
func (p Point) Round() Point {
	return Point{X: math.Floor(p.X + 0.5), Y: math.Floor(p.Y + 0.5), Z: math.Floor(p.Z + 0.5)}
}

// Scale multiplies the X,Y,Z components with a scalar value.
func (p Point) Scale(s float64) Point {
	return Point{X: p.X * s, Y: p.Y * s, Z: p.Z * s}
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

func (l Point) Cross(r Point) Point {
	return Point{X: l.Y*r.Z - l.Z*r.Y, Y: l.Z*r.X - l.X*r.Z, Z: l.X*r.Y - l.Y*r.X}
}

func (p Point) Mul(m Matrix) Point {
	x, y, z, _ := m.Transform(p.X, p.Y, p.Z, 1.0)
	return Point{X: x, Y: y, Z: z}
}
