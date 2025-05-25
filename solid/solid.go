package solid

import (
	"github.com/vibrantgio/seen"
	"github.com/vibrantgio/seen/face"
	"github.com/vibrantgio/seen/point"
	"github.com/vibrantgio/seen/transform"
)

type Solid struct {
	transform.Transform
	kind  string
	faces face.Faces
}

var _ seen.Object = (*Solid)(nil)

func NewSolid(kind string, csg CSG) seen.Object {
	s := Solid{transform.Default, kind, nil}
	for _, poly := range csg {
		points := make(point.Points, 0, len(poly.Vertices))
		for _, v := range poly.Vertices {
			points = append(points, point.Point{X: v.Pos.X, Y: v.Pos.Y, Z: v.Pos.Z})
		}
		s.faces = append(s.faces, face.FaceWith(points))
	}
	return &s
}

func (s Solid) Kind() string {
	return s.kind
}

func (s Solid) Faces() face.Faces {
	return s.faces
}
