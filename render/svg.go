package render

import (
	"github.com/reactivego/seen/document"
	"github.com/reactivego/seen"
	"github.com/reactivego/seen/affine"
)

// SvgRenderContext
type SvgRenderContext struct {
	layerAndContexts []svgLayerAndContext
	svg              *document.Element
}

func MakeSvgRenderContext(elementId string) RenderContext {
	c := &SvgRenderContext{}
	c.svg = document.GetElementById(elementId)
	return c
}

func (c *SvgRenderContext) Render() {
	c.Reset()
	for _, lc := range c.layerAndContexts {
		lc.context.Reset()
		lc.layer.Paint(lc.context)
		lc.context.Cleanup()
	}
	c.Cleanup()
}

func (c *SvgRenderContext) Animate() seen.RenderAnimator {
	return nil
}

func (c *SvgRenderContext) Layer(layer RenderLayer) {
	group := document.CreateElementNS(document.SVG_NS, "g")
	c.svg.AppendChild(group)
	lc := svgLayerAndContext{layer, MakeSvgPaintContext(group)}
	c.layerAndContexts = append(c.layerAndContexts, lc)
}

func (c *SvgRenderContext) Reset() {

}

func (c *SvgRenderContext) Cleanup() {

}

// svgLayerAndContext
type svgLayerAndContext struct {
	layer   RenderLayer
	context PaintContext
}

// SvgPaintContext
// implements PaintContext
type SvgPaintContext struct {
	group         *document.Element
	pathPainter   PathPainter
	textPainter   TextPainter
	circlePainter CirclePainter
	rectPainter   RectPainter
	i             int
}

func MakeSvgPaintContext(group *document.Element) PaintContext {
	c := &SvgPaintContext{}
	c.group = group
	c.pathPainter = MakeSvgPathPainter(c.elementFactory)
	c.textPainter = MakeSvgTextPainter(c.elementFactory)
	c.circlePainter = MakeSvgCirclePainter(c.elementFactory)
	c.rectPainter = MakeSvgRectPainter(c.elementFactory)
	return c
}

// Returns an element with tagname `type`.
//
// This method uses an internal iterator to add new elements as they are
// drawn. If there is no child element at the current index, we append one
// and return it. If an element exists at the current index and it is the
// same tag, we return that. If the element is a different type, we create
// one and replace it then return it.
func (c *SvgPaintContext) elementFactory(tag string) *document.Element {
	children := c.group.ChildNodes
	if c.i >= len(children) {
		path := document.CreateElementNS(document.SVG_NS, tag)
		c.group.AppendChild(path)
		c.i++
		return path
	}

	current := children[c.i]
	if current.Tag == tag {
		c.i++
		return current
	}

	path := document.CreateElementNS(document.SVG_NS, tag)
	c.group.ReplaceChild(path, current)
	c.i++
	return path
}

func (c *SvgPaintContext) Path() PathPainter {
	return c.pathPainter
}

func (c *SvgPaintContext) Rect() RectPainter {
	return c.rectPainter
}

func (c *SvgPaintContext) Circle() CirclePainter {
	return c.circlePainter
}

func (c *SvgPaintContext) Text() TextPainter {
	return c.textPainter
}

func (c *SvgPaintContext) Reset() {
	c.i = 0
}

func (c *SvgPaintContext) Cleanup() {
	children := c.group.ChildNodes
	for c.i < len(children) {
		children[c.i].SetAttribute("style", "display: none;")
		c.i++
	}
}

// SvgPainter
type SvgPainter struct {
	svgTag         string
	elementFactory func(tag string) *document.Element
	attributes     map[string]string
}

func MakeSvgPainter(svgTag string, elementFactory func(tag string) *document.Element) *SvgPainter {
	p := &SvgPainter{}
	if svgTag != "" {
		p.svgTag = svgTag
	} else {
		p.svgTag = "g"
	}
	p.elementFactory = elementFactory
	p.attributes = make(map[string]string)
	return p
}

func (p *SvgPainter) Clear() {
	p.attributes = make(map[string]string)
}

func (p *SvgPainter) Fill(style map[string]string) {
	p.Paint(style)
}

func (p *SvgPainter) Stroke(style map[string]string) {
	p.Paint(style)
}

func (p *SvgPainter) Paint(style map[string]string) {
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

// SvgPathPainter
type SvgPathPainter struct {
	*SvgPainter
}

func MakeSvgPathPainter(elementFactory func(tag string) *document.Element) *SvgPathPainter {
	return &SvgPathPainter{MakeSvgPainter("path", elementFactory)}
}

func (p *SvgPathPainter) Path(points []seen.Point) {
	str := "M"
	for _, point := range points {
		str += Fjoin(point.X, point.Y) + "L"
	}
	p.attributes["d"] = str[:len(str)-1]
}

// SvgTextPainter
type SvgTextPainter struct {
	*SvgPainter
}

func MakeSvgTextPainter(elementFactory func(tag string) *document.Element) *SvgTextPainter {
	return &SvgTextPainter{&SvgPainter{svgTag: "text", elementFactory: elementFactory}}
}

func (p *SvgTextPainter) FillText(t *affine.Matrix, text string, style map[string]string) {
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

// SvgCirclePainter
type SvgCirclePainter struct {
	*SvgPainter
}

func MakeSvgCirclePainter(elementFactory func(tag string) *document.Element) *SvgCirclePainter {
	return &SvgCirclePainter{MakeSvgPainter("circle", elementFactory)}
}

// SvgRectPainter
type SvgRectPainter struct {
	*SvgPainter
}

func MakeSvgRectPainter(elementFactory func(tag string) *document.Element) *SvgRectPainter {
	return &SvgRectPainter{MakeSvgPainter("rect", elementFactory)}
}

func (p *SvgRectPainter) Size(width, height float64) {
	p.attributes["width"] = Ftoa(width)
	p.attributes["height"] = Ftoa(height)
}

func (p *SvgRectPainter) CornerRadius(rx, ry float64) {
	p.attributes["rx"] = Ftoa(rx)
	p.attributes["ry"] = Ftoa(ry)
}
