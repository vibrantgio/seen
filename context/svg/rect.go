package svg

import "github.com/vibrantgio/seen/canvas"

// Rect
type Rect struct {
	*Styler
	precision *int
}

func newRect(elementFactory func(tag string) *Element, precision *int) *Rect {
	return &Rect{
		newStyler(func() *Element { return elementFactory("rect") }),
		precision,
	}
}

func (rect *Rect) Rect(width, height float64) canvas.RectPainter {
	rect.attributes["width"] = Ftoa(*rect.precision, width)
	rect.attributes["height"] = Ftoa(*rect.precision, height)
	return rect
}

func (rect *Rect) CornerRadius(rx, ry float64) canvas.RectPainter {
	rect.attributes["rx"] = Ftoa(*rect.precision, rx)
	rect.attributes["ry"] = Ftoa(*rect.precision, ry)
	return rect
}
