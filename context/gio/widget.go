package gio

import (
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
)

func Widget(tgt *Context, project func(w, h unit.Dp)) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		size := gtx.Constraints.Max
		project(gtx.Metric.PxToDp(size.X), gtx.Metric.PxToDp(size.Y))
		op.Affine(f32.NewAffine2D(gtx.Metric.PxPerDp, 0, 0, 0, gtx.Metric.PxPerDp, 0)).Add(gtx.Ops)
		tgt.Process(gtx)
		return layout.Dimensions{Size: size}
	}
}
