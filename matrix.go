package seen

import (
	"github.com/reactivego/seen/float"
	"github.com/reactivego/seen/transform"
)

// Matrix is a wrapper around transform.Matrix that can transform Point values.
type Matrix struct{ transform.Matrix }

var IdentityMatrix = Matrix{transform.IdentityMatrix}

func Translate(tx, ty, tz float64) Matrix {
	return Matrix{transform.Matrix{
		1, 0, 0, tx,
		0, 1, 0, ty,
		0, 0, 1, tz,
		0, 0, 0, 1,
	}}
}

func Scale(sx, sy, sz float64) Matrix {
	return Matrix{transform.Matrix{
		sx, 0, 0, 0,
		0, sy, 0, 0,
		0, 0, sz, 0,
		0, 0, 0, 1,
	}}
}

func (m Matrix) TransformPoint(p Point) Point {
	x, y, z, _ := m.Matrix.Transform(p.X, p.Y, p.Z, 1.0)
	return Point{x, y, z}
}

func (m Matrix) TransformCoordinate(p Coordinate) Coordinate {
	x, y, z, w := m.Matrix.Transform(p.X, p.Y, p.Z, p.W)
	return Coordinate{x, y, z, w}
}

func (l Matrix) Mul(r Matrix) Matrix {
	return Matrix{l.Matrix.Mul(r.Matrix)}
}

func (l Matrix) Equal(r Matrix) bool {
	for i, li := range l.Matrix {
		if !float.Equal(li, r.Matrix[i]) {
			return false
		}
	}
	return true
}

func (m Matrix) Scale(sx, sy, sz float64) Matrix {
	s := transform.Matrix{
		sx, 0, 0, 0,
		0, sy, 0, 0,
		0, 0, sz, 0,
		0, 0, 0, 1,
	}
	return Matrix{m.Matrix.Mul(s)}
}

func (m Matrix) Translate(tx, ty, tz float64) Matrix {
	s := transform.Matrix{
		1, 0, 0, tx,
		0, 1, 0, ty,
		0, 0, 1, tz,
		0, 0, 0, 1,
	}
	return Matrix{m.Matrix.Mul(s)}
}

func (m Matrix) TransformPoints(points []Point, transformedPoints []Point) (barycenter Point) {
	if len(transformedPoints) != len(points) {
		panic("internal error, slice lengths don't match")
	}

	// Create Barycenter point used in sorting surfaces in the painters algorithm
	barycenter = PointZero

	// Apply transform to points
	for i := range points {
		p := m.TransformPoint(points[i])
		transformedPoints[i] = p
		barycenter = barycenter.Add(p)
	}

	// Compute barycenter, which is used in sorting surfaces in the painters algorithm
	return barycenter.Scale(1.0 / float64(len(points)))
}

func (m Matrix) ProjectCoordinatesToPoints(coords []Coordinate, transformedPoints []Point) (barycenter Point) {
	if len(transformedPoints) != len(coords) {
		panic("internal error, slice lengths don't match")
	}

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
	return barycenter.Scale(1.0 / float64(len(coords)))
}
