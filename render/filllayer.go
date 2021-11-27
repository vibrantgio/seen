package render

import "github.com/reactivego/seen/color"

// FillLayer
// implements Layer
type FillLayer struct {
	Width, Height float64 //# 500,500
	Rx, Ry        float64
	Fill          string // fill: #EEE
}

func FillLayerWith(width, height, rx, ry float64, fill color.Color) *FillLayer {
	return &FillLayer{width, height, rx, ry, fill.Hex()}
}

func (l FillLayer) Paint(painter Painter) {
	rect := painter.Rect()
	rect.Size(l.Width, l.Height)
	if l.Rx != 0.0 || l.Ry != 0.0 {
		rect.CornerRadius(l.Rx, l.Ry)
	}
	rect.Fill(map[string]string{"fill": l.Fill})
}
