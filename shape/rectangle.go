package shape

import (
	"math"

	"github.com/reactivego/seen"
)

// Returns an axis-aligned 3D rectangle whose boundaries are defined by the
// two supplied points.
func Rectangle(point1, point2 seen.Point) *seen.Shape {
	compose := func(x, y, z func(float64, float64) float64) seen.Point {
		return seen.Point{
			X: x(point1.X, point2.X),
			Y: y(point1.Y, point2.Y),
			Z: z(point1.Z, point2.Z),
		}
	}
	points := []seen.Point{
		compose(math.Min, math.Min, math.Min),
		compose(math.Min, math.Min, math.Max),
		compose(math.Min, math.Max, math.Min),
		compose(math.Min, math.Max, math.Max),
		compose(math.Max, math.Min, math.Min),
		compose(math.Max, math.Min, math.Max),
		compose(math.Max, math.Max, math.Min),
		compose(math.Max, math.Max, math.Max),
	}
	return &seen.Shape{
		Type:      "rect",
		Transform: seen.DefaultTransform,
		Surfaces:  seen.SurfacesWith(points[:], CubeMap[:])}
}
