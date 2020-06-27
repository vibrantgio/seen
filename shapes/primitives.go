package shapes

import (
	"github.com/reactivego/seen"
	"github.com/reactivego/seen/affine"
)

// MakeCube returns a 2x2x2 cube, centered on the origin.
func MakeCube() *seen.Shape {
	points := [...]seen.Point{
	  {-1, -1, -1},
	  {-1, -1,  1},
	  {-1,  1, -1},
	  {-1,  1,  1},
	  { 1, -1, -1},
	  { 1, -1,  1},
	  { 1,  1, -1},
	  { 1,  1,  1},
	}
	s := &seen.Shape{}
	s.Init("cube", seen.MakeSurfaces(points[:],_CUBE_COORDINATE_MAP[:]))
	return s
}

// MakeUnitCube returns a 1x1x1 cube from the origin [0,0,0] to [1, 1, 1].
func MakeUnitCube() *seen.Shape {
	points := [...]seen.Point{
	  {0, 0, 0},
	  {0, 0, 1},
	  {0, 1, 0},
	  {0, 1, 1},
	  {1, 0, 0},
	  {1, 0, 1},
	  {1, 1, 0},
	  {1, 1, 1},
	}
	s := &seen.Shape{}
	s.Init("unitcube", seen.MakeSurfaces(points[:],_CUBE_COORDINATE_MAP[:]))
	return s
}

// MakeText creates a Shape for the given text.
// The text parameter is set as surface option "text" on the created shape.
// After the text is set the other surface options that were passed in via surfaceOptions
// are assigned to the surface options of the created shape.
//	font	e.g. "20px sans-serif" or "10px Roboto"
//	anchor	e.g. "middle"
func MakeText(text string, surfaceOptions map[string]string) *seen.Shape {
	surface := seen.MakeSurface(affine.ORTHONORMAL_BASIS)
	surface.Options["text"] = text
	for key, val := range surfaceOptions {
		surface.Options[key] = val
	}
	s := &seen.Shape{}
	s.Init("text", []seen.Surface{*surface})
	return s
}

// MakeIcosahedron returns an icosahedron that fits within a 2x2x2 cube, centered on the origin.
func MakeIcosahedron() *seen.Shape {
	s := &seen.Shape{}
	s.Init("icosahedron", seen.MakeSurfaces(_ICOSAHEDRON_POINTS[:], _ICOSAHEDRON_COORDINATE_MAP[:]))
	return s
}

// MakeSphere returns a sub-divided icosahedron, which approximates a sphere with
// triangles of equal size. A value of 2 is a good default for the subdivisions
// parameter, for every triangle of the original sphere 4*4= 16 triangles will
// be introduced. A subdivision value of 3 will generate 64 and a value
// of 4 will generated 256 triangles for every original triangle.
func MakeSphere(subdivisions int) *seen.Shape {

	triangles := make([][3]seen.Point,len(_ICOSAHEDRON_COORDINATE_MAP))
	for i,coords := range _ICOSAHEDRON_COORDINATE_MAP {
		for j,c := range coords {
			triangles[i][j] = _ICOSAHEDRON_POINTS[c]
		}
	}

	for i:=0; i<subdivisions; i++ {
		triangles = sphereSubdivideTriangles(triangles)
	}

	surfaces := make([]seen.Surface, len(triangles))
	for i,triangle := range triangles {
		surfaces[i] = *seen.MakeSurface(triangle[:])
	}

	s := &seen.Shape{}
	s.Init("sphere", surfaces)
	return s
}

// sphereSubdivideTriangles will return 4 triangles for every triangle passed in.
// Accepts an array of 3-tuples and returns an array of 3-tuples representing
// the triangular subdivision of the surface.
// The points of the triangle mesh passed in are supposed to be all unit
// vectors (length 1). Points introduced during triangulation are also normalized
// to unit length. The resulting mesh will therefore more and more start to approximate
// a unit sphere ater each call to sphereSubdivideTriangles.
func sphereSubdivideTriangles(triangles [][3]seen.Point) [][3]seen.Point {
	newTriangles := make([][3]seen.Point,0,len(triangles)*4)
	for _,tri := range triangles {
		v01 := tri[0].Add(&tri[1]).Normalize() // pull point back onto unit sphere.
		v12 := tri[1].Add(&tri[2]).Normalize()
		v20 := tri[2].Add(&tri[0]).Normalize()
		newTriangles = append(newTriangles, [3]seen.Point{tri[0], *v01, *v20})
		newTriangles = append(newTriangles, [3]seen.Point{tri[1], *v12, *v01})
		newTriangles = append(newTriangles, [3]seen.Point{tri[2], *v20, *v12})
		newTriangles = append(newTriangles, [3]seen.Point{*v01,   *v12, *v20})
	}
	return newTriangles
}

// Map to points in the surfaces of a cube
var _CUBE_COORDINATE_MAP = [...][]int{
  {0, 1, 3, 2}, // left
  {5, 4, 6, 7}, // right
  {1, 0, 4, 5}, // bottom
  {2, 3, 7, 6}, // top
  {3, 1, 5, 7}, // front
  {0, 2, 6, 4}, // back
}

const (
	_ICOX_X = 0.525731112119133606
	_ICOX_Z = 0.850650808352039932
)

var (
	_ICOSAHEDRON_COORDINATE_MAP = [...][]int{
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
	_ICOSAHEDRON_POINTS = [...]seen.Point{
		{-_ICOX_X, 0.0, -_ICOX_Z},
		{_ICOX_X, 0.0, -_ICOX_Z},
		{-_ICOX_X, 0.0, _ICOX_Z},
		{_ICOX_X, 0.0, _ICOX_Z},
		{0.0, _ICOX_Z, -_ICOX_X},
		{0.0, _ICOX_Z, _ICOX_X},
		{0.0, -_ICOX_Z, -_ICOX_X},
		{0.0, -_ICOX_Z, _ICOX_X},
		{_ICOX_Z, _ICOX_X, 0.0},
		{-_ICOX_Z, _ICOX_X, 0.0},
		{_ICOX_Z, -_ICOX_X, 0.0},
		{-_ICOX_Z, -_ICOX_X, 0.0},
	}
)
