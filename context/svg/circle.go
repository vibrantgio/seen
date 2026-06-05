package svg

import (
	"github.com/vibrantgio/seen/canvas"
	"github.com/vibrantgio/seen/point"
)

// Circle
type Circle struct {
	*Styler
	precision *int
}

func newCircle(elementFactory func(tag string) *Element, precision *int) *Circle {
	return &Circle{
		newStyler(func() *Element { return elementFactory("circle") }),
		precision,
	}
}

func (circle *Circle) Circle(center point.Point, radius float64) canvas.CirclePainter {
	circle.attributes["cx"] = Ftoa(*circle.precision, center.X)
	circle.attributes["cy"] = Ftoa(*circle.precision, center.Y)
	circle.attributes["r"] = Ftoa(*circle.precision, radius)
	return circle
}
