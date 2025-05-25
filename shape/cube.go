package shape

import (
	"math"

	"github.com/vibrantgio/seen"
	"github.com/vibrantgio/seen/face"
	"github.com/vibrantgio/seen/point"
)

// UnitCube returns a 1x1x1 cube from the origin [0,0,0] to [1, 1, 1].
func UnitCube() seen.Object {
	return NewShape("unitcube", UnitCubePoints[:], CubeFacets[:])
}

// UnitCubePoints form a 1x1x1 cube from the origin [0,0,0] to [1, 1, 1].
var UnitCubePoints = [...]point.Point{
	{X: 0, Y: 0, Z: 0},
	{X: 0, Y: 0, Z: 1},
	{X: 0, Y: 1, Z: 0},
	{X: 0, Y: 1, Z: 1},
	{X: 1, Y: 0, Z: 0},
	{X: 1, Y: 0, Z: 1},
	{X: 1, Y: 1, Z: 0},
	{X: 1, Y: 1, Z: 1},
}

// Returns an axis-aligned 3D rectangle whose boundaries are defined by the
// two supplied points.
func Rectangle(point1, point2 point.Point) seen.Object {
	compose := func(x, y, z func(float64, float64) float64) point.Point {
		return point.Point{
			X: x(point1.X, point2.X),
			Y: y(point1.Y, point2.Y),
			Z: z(point1.Z, point2.Z),
		}
	}
	points := [...]point.Point{
		compose(math.Min, math.Min, math.Min),
		compose(math.Min, math.Min, math.Max),
		compose(math.Min, math.Max, math.Min),
		compose(math.Min, math.Max, math.Max),
		compose(math.Max, math.Min, math.Min),
		compose(math.Max, math.Min, math.Max),
		compose(math.Max, math.Max, math.Min),
		compose(math.Max, math.Max, math.Max),
	}
	return NewShape("rect", points[:], CubeFacets[:])
}

// Cube returns a 2x2x2 cube, centered on the origin.
func Cube() seen.Object {
	return NewShape("cube", CubePoints[:], CubeFacets[:])
}

// CubePoints form a 2x2x2 cube, centered on the origin.
var CubePoints = [...]point.Point{
	{X: -1, Y: -1, Z: -1},
	{X: -1, Y: -1, Z: 1},
	{X: -1, Y: 1, Z: -1},
	{X: -1, Y: 1, Z: 1},
	{X: 1, Y: -1, Z: -1},
	{X: 1, Y: -1, Z: 1},
	{X: 1, Y: 1, Z: -1},
	{X: 1, Y: 1, Z: 1},
}

// Map to points in the faces of a cube
var CubeFacets = [...]face.Facet{
	{0, 1, 3, 2}, // left
	{5, 4, 6, 7}, // right
	{1, 0, 4, 5}, // bottom
	{2, 3, 7, 6}, // top
	{3, 1, 5, 7}, // front
	{0, 2, 6, 4}, // back
}
