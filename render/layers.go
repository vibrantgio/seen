package render

// RenderLayer
type RenderLayer interface {
	Paint(context PaintContext)
}

// FillLayer
// implements RenderLayer
type FillLayer struct {
	Width, Height float64 //# 500,500
	Rx, Ry        float64
	Fill          string  // fill: #EEE
}

func MakeFillLayer(width, height, rx, ry float64, fill string) *FillLayer {
	return &FillLayer{width, height, rx, ry, fill}
}

func (l *FillLayer) Paint(context PaintContext) {
	rectPainter := context.Rect()
	rectPainter.Size(l.Width, l.Height)
	if l.Rx != 0.0 || l.Ry != 0.0 {
		rectPainter.CornerRadius(l.Rx,l.Ry)
	}
	rectPainter.Fill(map[string]string{"fill": l.Fill})
}
