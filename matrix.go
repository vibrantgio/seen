package seen

import (
	"github.com/reactivego/seen/float"
	"github.com/reactivego/seen/transform"
)

// Matrix is a wrapper around transform.Mat4x4 that can transform Point values.
type Matrix struct{ transform.Mat4x4 }

var IdentityMatrix = Matrix{transform.IdentityMat4x4}

func Translate(tx, ty, tz float64) Matrix {
	return Matrix{transform.Mat4x4{
		1, 0, 0, tx,
		0, 1, 0, ty,
		0, 0, 1, tz,
		0, 0, 0, 1,
	}}
}

func Scale(sx, sy, sz float64) Matrix {
	return Matrix{transform.Mat4x4{
		sx, 0, 0, 0,
		0, sy, 0, 0,
		0, 0, sz, 0,
		0, 0, 0, 1,
	}}
}

func (m Matrix) TransformPoint(p Point) Point {
	x, y, z, _ := m.Mat4x4.Transform(p.X, p.Y, p.Z, 1.0)
	return Point{x, y, z}
}

func (m Matrix) TransformCoordinate(p Coordinate) Coordinate {
	x, y, z, w := m.Mat4x4.Transform(p.X, p.Y, p.Z, p.W)
	return Coordinate{x, y, z, w}
}

func (l Matrix) Mul(r Matrix) Matrix {
	return Matrix{l.Mat4x4.Mul(r.Mat4x4)}
}

func (l Matrix) Equal(r Matrix) bool {
	for i, li := range l.Mat4x4 {
		if !float.Equal(li, r.Mat4x4[i]) {
			return false
		}
	}
	return true
}

func (m Matrix) Scale(sx, sy, sz float64) Matrix {
	s := transform.Mat4x4{
		sx, 0, 0, 0,
		0, sy, 0, 0,
		0, 0, sz, 0,
		0, 0, 0, 1,
	}
	return Matrix{m.Mat4x4.Mul(s)}
}

func (m Matrix) Translate(tx, ty, tz float64) Matrix {
	s := transform.Mat4x4{
		1, 0, 0, tx,
		0, 1, 0, ty,
		0, 0, 1, tz,
		0, 0, 0, 1,
	}
	return Matrix{m.Mat4x4.Mul(s)}
}

func (m Matrix) TransformPoints(points []Point) (transformedPoints []Point, barycenter Point) {
	// Length of the passed in points slice
	pointsLen := len(points)

	// Size the Points array to fit the length of the points slice passed in.
	transformedPoints = make([]Point, pointsLen)

	// Create Barycenter point used in sorting surfaces in the painters algorithm
	barycenter = PointZero

	// Apply transform to points
	for i := range points {
		p := m.TransformPoint(points[i])
		transformedPoints[i] = p
		barycenter = barycenter.Add(p)
	}

	// Compute barycenter, which is used in sorting surfaces in the painters algorithm
	barycenter = barycenter.Scale(1.0 / float64(pointsLen))
	return transformedPoints, barycenter
}

func (m Matrix) ProjectCoordinatesToPoints(coords []Coordinate) (transformedPoints []Point, barycenter Point) {

	// Length of the passed in coords slice
	coordsLen := len(coords)

	// Size the Coordinates array to fit the length of the coords slice passed in.
	transformedPoints = make([]Point, coordsLen)

	barycenter = PointZero

	// Apply transform to coords
	for i := range coords {
		// Transform the homogeneous coordinate
		c := m.TransformCoordinate(coords[i])

		// Calling ToPoint on a Coordinate will apply the clip so it scales the x and y
		// coordinates in a perspective projection. This is done by dividing the X,Y,Z
		// components by c.W
		p := c.ToPoint()

		// copy p into the points array.
		transformedPoints[i] = p
		barycenter = barycenter.Add(p)
	}

	// Compute barycenter, which is used in sorting surfaces in the painters algorithm
	barycenter = barycenter.Scale(1.0 / float64(coordsLen))

	return transformedPoints, barycenter
}
