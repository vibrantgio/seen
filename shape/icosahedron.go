package shape

import "github.com/reactivego/seen"

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

// Icosahedron returns an icosahedron that fits within a 2x2x2 cube, centered on the origin.
func Icosahedron() *seen.Shape {
	return &seen.Shape{
		Type:      "icosahedron",
		Transform: seen.DefaultTransform,
		Surfaces:  seen.SurfacesWith(IcosahedronPoints[:], IcosahedronMap[:])}
}

// Sphere returns a sub-divided icosahedron, which approximates a sphere with
// triangles of equal size. A value of 2 is a good default for the subdivisions
// parameter, for every triangle of the original sphere 4*4= 16 triangles will
// be introduced. A subdivision value of 3 will generate 64 and a value
// of 4 will generated 256 triangles for every original triangle.
func Sphere(subdivisions int) *seen.Shape {

	triangles := make([][3]seen.Point, len(IcosahedronMap))
	for i, coords := range IcosahedronMap {
		for j, c := range coords {
			triangles[i][j] = IcosahedronPoints[c]
		}
	}

	// Subdivide the icosahedron. Every subdivision returns 4 triangles for every triangle
	for i := 0; i < subdivisions; i++ {
		/// The mesh will approximate a unit sphere more with every subdivide.
		subdivide := func(triangles [][3]seen.Point) [][3]seen.Point {
			newTriangles := make([][3]seen.Point, 0, len(triangles)*4)
			// The points of the triangle mesh passed in are supposed to be all unit
			// vectors (length 1).
			for _, tri := range triangles {
				// Points introduced during triangulation are also normalized to unit length.
				v01 := tri[0].Add(tri[1]).Normalize() // pull point back onto unit sphere.
				v12 := tri[1].Add(tri[2]).Normalize()
				v20 := tri[2].Add(tri[0]).Normalize()
				newTriangles = append(newTriangles, [3]seen.Point{tri[0], v01, v20})
				newTriangles = append(newTriangles, [3]seen.Point{tri[1], v12, v01})
				newTriangles = append(newTriangles, [3]seen.Point{tri[2], v20, v12})
				newTriangles = append(newTriangles, [3]seen.Point{v01, v12, v20})
			}
			return newTriangles
		}
		triangles = subdivide(triangles)
	}

	surfaces := make(seen.Surfaces, len(triangles))
	for i, triangle := range triangles {
		surfaces[i] = *seen.SurfaceWith(triangle[:])
	}

	return &seen.Shape{Type: "sphere", Transform: seen.DefaultTransform, Surfaces: surfaces}
}
