package svg

import (
	"github.com/vibrantgio/seen/canvas"
	"github.com/vibrantgio/seen/point"
)

// Circle
type Circle struct{ *Styler }

func newCircle(elementFactory func(tag string) *Element) *Circle {
	return &Circle{
		newStyler(func() *Element { return elementFactory("circle") }),
	}
}

func (circle *Circle) Circle(center point.Point, radius float64) canvas.CirclePainter {
	circle.attributes["cx"] = Ftoa(center.X)
	circle.attributes["cy"] = Ftoa(center.Y)
	circle.attributes["r"] = Ftoa(radius)
	return circle
}
