package bsort

import (
	"github.com/vibrantgio/seen"
	"github.com/vibrantgio/seen/layer/bsort/bsp"
	"github.com/vibrantgio/seen/matrix"
)

type Planes []bsp.Plane

var _ seen.Handler = (*Planes)(nil)

func (b Planes) EnterGroup() {}

func (b Planes) LeaveGroup() {}

func (b Planes) VisitLight(light seen.Light, model matrix.Matrix) {}

func (b *Planes) VisitObject(object seen.Object, model matrix.Matrix) {
	for _, face := range object.Faces() {
		*b = append(*b, bsp.PlaneWith(face.Id, face.Points, model))
	}
}
