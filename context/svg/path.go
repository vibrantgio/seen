package svg

import (
	"github.com/vibrantgio/seen/canvas"
	"github.com/vibrantgio/seen/point"
)

// Path
type Path struct {
	*Styler
	precision *int
}

func newPath(elementFactory func(tag string) *Element, precision *int) *Path {
	return &Path{
		newStyler(func() *Element { return elementFactory("path") }),
		precision,
	}
}

func (p *Path) Path(points []point.Point) canvas.PathPainter {
	str := "M"
	for _, point := range points {
		str += Fjoin(*p.precision, point.X, point.Y) + "L"
	}
	p.attributes["d"] = str[:len(str)-1]
	return p
}
