package svg

import (
	"github.com/vibrantgio/seen/affine"
	"github.com/vibrantgio/seen/canvas"
)

// Text
type Text struct {
	makeElement func() *Element
}

func newText(elementFactory func(tag string) *Element) *Text {
	return &Text{func() *Element {
		return elementFactory("text")
	}}
}

func (p *Text) FillText(t affine.Matrix, text string, style canvas.Style) {
	el := p.makeElement()

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
