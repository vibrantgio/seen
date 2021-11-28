package shape

import (
	"github.com/reactivego/seen"
)

// UnitCubePoints form a 1x1x1 cube from the origin [0,0,0] to [1, 1, 1].
var UnitCubePoints = [...]seen.Point{
	{X: 0, Y: 0, Z: 0},
	{X: 0, Y: 0, Z: 1},
	{X: 0, Y: 1, Z: 0},
	{X: 0, Y: 1, Z: 1},
	{X: 1, Y: 0, Z: 0},
	{X: 1, Y: 0, Z: 1},
	{X: 1, Y: 1, Z: 0},
	{X: 1, Y: 1, Z: 1},
}

// UnitCube returns a 1x1x1 cube from the origin [0,0,0] to [1, 1, 1].
func UnitCube() *seen.Shape {
	return &seen.Shape{
		Type:      "unitcube",
		Transform: seen.DefaultTransform,
		Surfaces:  seen.SurfacesWith(UnitCubePoints[:], CubeMap[:])}
}
