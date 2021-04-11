package shape

import (
	"math"

	"github.com/reactivego/seen"
	"github.com/reactivego/seen/affine"
)

// Cube returns a 2x2x2 cube, centered on the origin.
func Cube() *seen.Shape {
	s := seen.ShapeWith("cube", seen.SurfacesWith(CubePoints[:], CubeCoordinateMap[:]))
	return &s
}

// UnitCube returns a 1x1x1 cube from the origin [0,0,0] to [1, 1, 1].
func UnitCube() *seen.Shape {
	s := seen.ShapeWith("unitcube", seen.SurfacesWith(UnitCubePoints[:], CubeCoordinateMap[:]))
	return &s
}

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

	s := seen.ShapeWith("rect", seen.SurfacesWith(
		points[:], CubeCoordinateMap[:]))
	return &s
}

// Text creates a Shape for the given text.
// The text parameter is set as surface option "text" on the created shape.
// After the text is set the other surface options that were passed in via surfaceOptions
// are assigned to the surface options of the created shape.
//	font	e.g. "20px sans-serif" or "10px Roboto"
//	anchor	e.g. "middle"
func Text(text string, surfaceOptions map[string]string) *seen.Shape {
	surface := seen.SurfaceWith(affine.ORTHONORMAL_BASIS)
	surface.Options["text"] = text
	for key, val := range surfaceOptions {
		surface.Options[key] = val
	}
	s := seen.ShapeWith("text",
		[]seen.Surface{*surface})
	return &s
}

// Icosahedron returns an icosahedron that fits within a 2x2x2 cube, centered on the origin.
func Icosahedron() *seen.Shape {
	s := seen.ShapeWith("icosahedron", seen.SurfacesWith(
		IcosahedronPoints[:],
		IcosahedronCoordinateMap[:]),
	)
	return &s
}

// Sphere returns a sub-divided icosahedron, which approximates a sphere with
// triangles of equal size. A value of 2 is a good default for the subdivisions
// parameter, for every triangle of the original sphere 4*4= 16 triangles will
// be introduced. A subdivision value of 3 will generate 64 and a value
// of 4 will generated 256 triangles for every original triangle.
func Sphere(subdivisions int) *seen.Shape {

	triangles := make([][3]seen.Point, len(IcosahedronCoordinateMap))
	for i, coords := range IcosahedronCoordinateMap {
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

	surfaces := make([]seen.Surface, len(triangles))
	for i, triangle := range triangles {
		surfaces[i] = *seen.SurfaceWith(triangle[:])
	}

	s := seen.ShapeWith("sphere", surfaces)
	return &s
}
