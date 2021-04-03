package seen

import (
	"math"

	"github.com/reactivego/seen/float"
)

// Point is an object used for specifying the outline of a Surface.
type Point struct {
	X, Y, Z float64
}

// PointZero is the zero point
var PointZero = Point{}

// PointX is the unit X axis.
var PointX = Point{1, 0, 0}

// PointY is the unit Y axis.
var PointY = Point{0, 1, 0}

// PointZ is the unit Z axis.
var PointZ = Point{0, 0, 1}

// MakePointNormal returns the normal for a slice of points.
// Can be used for e.g. backface culling or shading.
func PointNormal(points []Point) Point {
	pointsLen := len(points)
	if pointsLen < 2 {
		return PointZ // Default normal
	}
	p0 := points[0]
	p1 := points[1]
	p2 := points[pointsLen-1]
	v0 := p1.Subtract(p0)
	v1 := p2.Subtract(p0)
	return v0.Cross(v1).Normalize()
}

// Normalize returns a pointer to a copy of the Point after
// normalizing the copy so the vector length is 1.0
func (p Point) Normalize() Point {
	mag := p.Length()
	return Point{p.X / mag, p.Y / mag, p.Z / mag}
}

// Subtract returns a pointer to a copy of this Point after
// subtracting the passed in Point r from this copy.
// The W values are not subtracted.
func (l Point) Subtract(r Point) Point {
	return Point{l.X - r.X, l.Y - r.Y, l.Z - r.Z}
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
	return Point{l.X + r.X, l.Y + r.Y, l.Z + r.Z}
}

// Round rounds the X,Y,Z components to the nearest integer value.
func (p Point) Round() Point {
	return Point{math.Floor(p.X + 0.5), math.Floor(p.Y + 0.5), math.Floor(p.Z + 0.5)}
}

// Scale multiplies the X,Y,Z components with a scalar value.
func (p Point) Scale(s float64) Point {
	return Point{p.X * s, p.Y * s, p.Z * s}
}

func (l Point) Cross(r Point) Point {
	return Point{l.Y*r.Z - l.Z*r.Y, l.Z*r.X - l.X*r.Z, l.X*r.Y - l.Y*r.X}
}

func (l Point) Dot(r Point) float64 {
	return l.X*r.X + l.Y*r.Y + l.Z*r.Z
}

func (p Point) ToCoordinate() Coordinate {
	return Coordinate{p.X, p.Y, p.Z, 1.0}
}

type Coordinate struct {
	X, Y, Z, W float64
}

func (hc Coordinate) ToPoint() Point {
	return Point{hc.X / hc.W, hc.Y / hc.W, hc.Z / hc.W}
}
