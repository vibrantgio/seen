package svg

import (
	"github.com/reactivego/seen/affine"
	"github.com/reactivego/seen/document"
	"github.com/reactivego/seen/render"
)

// TextPainter
type TextPainter struct {
	*CommonPainter
}

func NewTextPainter(elementFactory func(tag string) *document.Element) *TextPainter {
	return &TextPainter{&CommonPainter{svgTag: "text", elementFactory: elementFactory}}
}

func (p *TextPainter) FillText(t affine.Matrix, text string, style render.Style) {
	el := p.elementFactory(p.svgTag)

	// set the transform attribute given the matrix m
	el.SetAttribute("transform", "matrix("+Fjoin(t.A, t.B, t.C, t.D, t.E, t.F)+")")

	// serialize the style map.
	str := ""
	for key, value := range style {
		if value != "" {
			str += key + ":" + value + ";"
		}
	}
	el.SetAttribute("style", str)
	el.TextContent = text
}
