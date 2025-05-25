package shape

import (
	"github.com/vibrantgio/seen"
	"github.com/vibrantgio/seen/face"
	"github.com/vibrantgio/seen/point"
)

// Sphere returns a sub-divided icosahedron, which approximates a sphere with
// triangles of equal size. With a subdivisions value of 2, the 20 triangles
// that form the icosahedron are split into 4 triangles twice, creating
// 20 x 4 x 4 = 320 triangles.
func Sphere(subdivisions int) seen.Object {

	triangles := make([][3]point.Point, len(IcosahedronFacets))
	for i, facet := range IcosahedronFacets {
		triangle := &triangles[i]
		for j, index := range facet {
			triangle[j] = IcosahedronPoints[index]
		}
	}

	// Subdivide the icosahedron. Every subdivision returns 4 triangles for every triangle
	for range subdivisions {
		/// The mesh will approximate a unit sphere more with every subdivide.
		subdivide := func(triangles [][3]point.Point) [][3]point.Point {
			newTriangles := make([][3]point.Point, 0, len(triangles)*4)
			// The points of the triangle mesh passed in are supposed to be all unit
			// vectors (length 1).
			for _, tri := range triangles {
				// Points introduced during triangulation are also normalized to unit length.
				v01 := tri[0].Plus(tri[1]).Normalize() // pull point back onto unit sphere.
				v12 := tri[1].Plus(tri[2]).Normalize()
				v20 := tri[2].Plus(tri[0]).Normalize()
				newTriangles = append(newTriangles, [3]point.Point{tri[0], v01, v20})
				newTriangles = append(newTriangles, [3]point.Point{tri[1], v12, v01})
				newTriangles = append(newTriangles, [3]point.Point{tri[2], v20, v12})
				newTriangles = append(newTriangles, [3]point.Point{v01, v12, v20})
			}
			return newTriangles
		}
		triangles = subdivide(triangles)
	}

	faces := make(face.Faces, 0, len(triangles))
	for _, triangle := range triangles {
		faces = append(faces, face.FaceWith(triangle[:]))
	}

	return NewShapeWithFaces("sphere", faces)
}

// Icosahedron returns an icosahedron that fits within a 2x2x2 cube, centered on the origin.
func Icosahedron() seen.Object {
	return NewShape("icosahedron", IcosahedronPoints[:], IcosahedronFacets[:])
}

const IcoXX = 0.525731112119133606
const IcoXZ = 0.850650808352039932

var IcosahedronPoints = [...]point.Point{
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

var IcosahedronFacets = [...]face.Facet{
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
