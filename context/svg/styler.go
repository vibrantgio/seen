package svg

import "github.com/vibrantgio/seen/canvas"

// Styler
type Styler struct {
	makeElement func() *Element
	attributes  map[string]string
}

func newStyler(makeElement func() *Element) *Styler {
	return &Styler{
		makeElement: makeElement,
		attributes:  make(map[string]string),
	}
}

func (p *Styler) Clear() {
	clear(p.attributes)
}

func (p *Styler) Fill(style canvas.Style) {
	p.Paint(style)
}

func (p *Styler) Stroke(style canvas.Style) {
	p.Paint(style)
}

func (p *Styler) Paint(style canvas.Style) {
	el := p.makeElement()

	str := ""
	for key, value := range style {
		str += key + ":" + value + ";"
	}
	el.SetAttribute("style", str)

	for key, value := range p.attributes {
		el.SetAttribute(key, value)
	}
}
