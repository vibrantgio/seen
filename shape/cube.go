package shape

import "github.com/reactivego/seen"

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

// Map to points in the surfaces of a cube
var CubeCoordinateMap = [...][]int{
	{0, 1, 3, 2}, // left
	{5, 4, 6, 7}, // right
	{1, 0, 4, 5}, // bottom
	{2, 3, 7, 6}, // top
	{3, 1, 5, 7}, // front
	{0, 2, 6, 4}, // back
}
