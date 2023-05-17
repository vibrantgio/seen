package svg

import (
	"github.com/reactivego/seen"
	"github.com/reactivego/seen/document"
)

// PathPainter
type PathPainter struct {
	*CommonPainter
}

func NewPathPainter(elementFactory func(tag string) *document.Element) *PathPainter {
	return &PathPainter{NewCommonPainter("path", elementFactory)}
}

func (p *PathPainter) Path(points []seen.Point) {
	str := "M"
	for _, point := range points {
		str += Fjoin(point.X, point.Y) + "L"
	}
	p.attributes["d"] = str[:len(str)-1]
}
