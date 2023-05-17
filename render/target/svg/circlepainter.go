package svg

import "github.com/reactivego/seen/document"

// CirclePainter
type CirclePainter struct {
	*CommonPainter
}

func NewCirclePainter(elementFactory func(tag string) *document.Element) *CirclePainter {
	return &CirclePainter{NewCommonPainter("circle", elementFactory)}
}
