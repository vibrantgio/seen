package backdrop

import (
	"github.com/vibrantgio/seen/canvas"
	"github.com/vibrantgio/seen/color"
	"github.com/vibrantgio/seen/layer"
)

// Layer implements layer.Layer
type Layer struct {
	Width, Height float64 //# 500,500
	Rx, Ry        float64
	Fill          string // fill: #EEE
}

var _ layer.Layer = (*Layer)(nil)

func NewLayer(width, height, rx, ry float64, fill color.Color) *Layer {
	return &Layer{width, height, rx, ry, fill.Hex()}
}

func (l Layer) RenderOn(c canvas.Canvas) {
	rect := c.Rect().Rect(l.Width, l.Height)
	if l.Rx != 0.0 || l.Ry != 0.0 {
		rect = rect.CornerRadius(l.Rx, l.Ry)
	}
	rect.Fill(canvas.Style{"fill": l.Fill})
}
