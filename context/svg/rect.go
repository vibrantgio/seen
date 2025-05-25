package svg

import "github.com/vibrantgio/seen/canvas"

// Rect
type Rect struct {
	*Styler
}

func newRect(elementFactory func(tag string) *Element) *Rect {
	return &Rect{
		newStyler(func() *Element { return elementFactory("rect") }),
	}
}

func (rect *Rect) Rect(width, height float64) canvas.RectPainter {
	rect.attributes["width"] = Ftoa(width)
	rect.attributes["height"] = Ftoa(height)
	return rect
}

func (rect *Rect) CornerRadius(rx, ry float64) canvas.RectPainter {
	rect.attributes["rx"] = Ftoa(rx)
	rect.attributes["ry"] = Ftoa(ry)
	return rect
}
