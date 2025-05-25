package shape

import (
	"github.com/vibrantgio/seen"
	"github.com/vibrantgio/seen/face"
	"github.com/vibrantgio/seen/point"
	"github.com/vibrantgio/seen/transform"
)

// A shape like a polygon is made up of vertices, edges, and faces, which define
// its structure in 3D space. Vertices are points in space, edges connect these
// points, and faces are the flat faces enclosed by edges, forming the
// overall geometry of the model. They may create a closed 3D shape, but not
// necessarily. For example, a cube is a closed shape, but a patch is not.
type shape struct {
	transform.Transform
	kind  string
	faces face.Faces
}

var _ seen.Object = (*shape)(nil)

func NewShape(kind string, points point.Points, facets face.Facets) seen.Object {
	return NewShapeWithFaces(kind, facets.FacesWith(points))
}

func NewShapeWithFaces(kind string, faces face.Faces) seen.Object {
	return &shape{transform.Default, kind, faces}
}

func (shape shape) Kind() string {
	return shape.kind
}

func (shape shape) Faces() face.Faces {
	return shape.faces
}
