package svg

import "github.com/reactivego/seen/document"

// CommonPainter
type CommonPainter struct {
	svgTag         string
	elementFactory func(tag string) *document.Element
	attributes     map[string]string
}

func NewCommonPainter(svgTag string, elementFactory func(tag string) *document.Element) *CommonPainter {
	p := &CommonPainter{}
	if svgTag != "" {
		p.svgTag = svgTag
	} else {
		p.svgTag = "g"
	}
	p.elementFactory = elementFactory
	p.attributes = make(map[string]string)
	return p
}

func (p *CommonPainter) Clear() {
	p.attributes = make(map[string]string)
}

func (p *CommonPainter) Fill(style map[string]string) {
	p.Paint(style)
}

func (p *CommonPainter) Stroke(style map[string]string) {
	p.Paint(style)
}

func (p *CommonPainter) Paint(style map[string]string) {
	el := p.elementFactory(p.svgTag)

	str := ""
	for key, value := range style {
		str += key + ":" + value + ";"
	}
	el.SetAttribute("style", str)

	for key, value := range p.attributes {
		el.SetAttribute(key, value)
	}
}
