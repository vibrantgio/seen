package shape

import "github.com/reactivego/seen"

// Icosahedron returns an icosahedron that fits within a 2x2x2 cube, centered on the origin.
func Icosahedron() *seen.Shape {
	return &seen.Shape{
		Type:      "icosahedron",
		Transform: seen.DefaultTransform,
		Surfaces:  seen.SurfacesWith(IcosahedronPoints[:], IcosahedronMap[:])}
}

const IcoXX = 0.525731112119133606
const IcoXZ = 0.850650808352039932

var IcosahedronPoints = [...]seen.Point{
	{X: -IcoXX, Y: 0.0, Z: -IcoXZ},
	{X: IcoXX, Y: 0.0, Z: -IcoXZ},
	{X: -IcoXX, Y: 0.0, Z: IcoXZ},
	{X: IcoXX, Y: 0.0, Z: IcoXZ},
	{X: 0.0, Y: IcoXZ, Z: -IcoXX},
	{X: 0.0, Y: IcoXZ, Z: IcoXX},
	{X: 0.0, Y: -IcoXZ, Z: -IcoXX},
	{X: 0.0, Y: -IcoXZ, Z: IcoXX},
	{X: IcoXZ, Y: IcoXX, Z: 0.0},
	{X: -IcoXZ, Y: IcoXX, Z: 0.0},
	{X: IcoXZ, Y: -IcoXX, Z: 0.0},
	{X: -IcoXZ, Y: -IcoXX, Z: 0.0},
}

var IcosahedronMap = [...][]int{
	{0, 4, 1},
	{0, 9, 4},
	{9, 5, 4},
	{4, 5, 8},
	{4, 8, 1},
	{8, 10, 1},
	{8, 3, 10},
	{5, 3, 8},
	{5, 2, 3},
	{2, 7, 3},
	{7, 10, 3},
	{7, 6, 10},
	{7, 11, 6},
	{11, 0, 6},
	{0, 1, 6},
	{6, 1, 10},
	{9, 0, 11},
	{9, 11, 2},
	{9, 2, 5},
	{7, 2, 11},
}
