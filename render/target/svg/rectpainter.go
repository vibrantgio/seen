package svg

import "github.com/reactivego/seen/document"

// RectPainter
type RectPainter struct {
	*CommonPainter
}

func NewRectPainter(elementFactory func(tag string) *document.Element) *RectPainter {
	return &RectPainter{NewCommonPainter("rect", elementFactory)}
}

func (p *RectPainter) Size(width, height float64) {
	p.attributes["width"] = Ftoa(width)
	p.attributes["height"] = Ftoa(height)
}

func (p *RectPainter) CornerRadius(rx, ry float64) {
	p.attributes["rx"] = Ftoa(rx)
	p.attributes["ry"] = Ftoa(ry)
}
