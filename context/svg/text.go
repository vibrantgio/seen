package svg

import (
	"github.com/vibrantgio/seen/affine"
	"github.com/vibrantgio/seen/canvas"
)

// Text
type Text struct {
	makeElement func() *Element
	precision   *int
}

func newText(elementFactory func(tag string) *Element, precision *int) *Text {
	return &Text{
		makeElement: func() *Element { return elementFactory("text") },
		precision:   precision,
	}
}

func (p *Text) FillText(t affine.Matrix, text string, style canvas.Style) {
	el := p.makeElement()

	// set the transform attribute given the matrix m
	el.SetAttribute("transform", "matrix("+Fjoin(*p.precision, t.A, t.B, t.C, t.D, t.E, t.F)+")")

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
