package shape

import "github.com/reactivego/seen"

// Sphere returns a sub-divided icosahedron, which approximates a sphere with
// triangles of equal size. With a subdivisions value of 2, the 20 triangles
// that form the icosahedron are split into 4 triangles twice, creating
// 20 x 4 x 4 = 320 triangles.
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
