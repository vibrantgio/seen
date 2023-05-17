package shape

import (
	"github.com/reactivego/seen"
)

// Cube returns a 2x2x2 cube, centered on the origin.
func Cube() *seen.Shape {
	return &seen.Shape{
		Type:      "cube",
		Transform: seen.DefaultTransform,
		Surfaces:  seen.SurfacesWith(CubePoints[:], CubeMap[:])}
}

// CubePoints form a 2x2x2 cube, centered on the origin.
var CubePoints = [...]seen.Point{
	{X: -1, Y: -1, Z: -1},
	{X: -1, Y: -1, Z: 1},
	{X: -1, Y: 1, Z: -1},
	{X: -1, Y: 1, Z: 1},
	{X: 1, Y: -1, Z: -1},
	{X: 1, Y: -1, Z: 1},
	{X: 1, Y: 1, Z: -1},
	{X: 1, Y: 1, Z: 1},
}

// Map to points in the surfaces of a cube
var CubeMap = [...][]int{
	{0, 1, 3, 2}, // left
	{5, 4, 6, 7}, // right
	{1, 0, 4, 5}, // bottom
	{2, 3, 7, 6}, // top
	{3, 1, 5, 7}, // front
	{0, 2, 6, 4}, // back
}
