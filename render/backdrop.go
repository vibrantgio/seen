package render

import "github.com/reactivego/seen/color"

// Backdrop implements render.Layer
type Backdrop struct {
	Width, Height float64 //# 500,500
	Rx, Ry        float64
	Fill          string // fill: #EEE
}

func NewBackdrop(width, height, rx, ry float64, fill color.Color) *Backdrop {
	return &Backdrop{width, height, rx, ry, fill.Hex()}
}

func (l Backdrop) Paint(painter Painter) {
	rect := painter.Rect()
	rect.Size(l.Width, l.Height)
	if l.Rx != 0.0 || l.Ry != 0.0 {
		rect.CornerRadius(l.Rx, l.Ry)
	}
	rect.Fill(map[string]string{"fill": l.Fill})
}
