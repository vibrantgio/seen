package bsort

import (
	"math"

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
		plane := bsp.PlaneWith(face.Id, face.Points, model)
		if math.IsNaN(plane.Normal.X) {
			// Degenerate face (collinear or duplicated points): it spans no
			// plane and has no projected area, so skip it rather than let a
			// NaN normal corrupt the BSP side classification.
			continue
		}
		if _, text := face.Options["text"]; text {
			// A text face paints its whole string from its points; pieces
			// of it would each repeat the text, so it must never be cut.
			plane.NoSplit = true
		}
		*b = append(*b, plane)
	}
}
