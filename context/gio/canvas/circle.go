package canvas

import (
	"gioui.org/op"
	"github.com/vibrantgio/seen/canvas"
	"github.com/vibrantgio/seen/point"
)

// Circle
type Circle struct {
	*op.Ops
	center point.Point
	radius float64
}

func (circle *Circle) Circle(center point.Point, radius float64) canvas.CirclePainter {
	circle.center = center
	circle.radius = radius
	return circle
}

func (circle *Circle) Fill(style canvas.Style) {
}
