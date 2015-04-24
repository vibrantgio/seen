package seen

import (
	"strconv"
	"strings"
	"xpt.nl/document"
)

const SVG_NS = "http://www.w3.org/2000/svg"

func formatFloat(v float64) string {
	return strconv.FormatFloat(v, 'f', -1, 64)
}

func joinFloats(v ...float64) string {
	if len(v) == 0 {
		return ""
	}
	s := []byte(formatFloat(v[0]))
	for _, f := range v[1:] {
		s = append(append(s, ' '), formatFloat(f)...)
	}
	return string(s)
}

//----------------------------
// RenderContext
//----------------------------

type RenderContext interface {
	Render()
	Animate() RenderAnimator
	Layer(layer RenderLayer)

	Reset()
	Cleanup()
}

//----------------------------
// RenderAnimator
//----------------------------

type RenderAnimator interface {
}

//----------------------------
// RenderLayer
//----------------------------

type RenderLayer interface {
	Render(context RenderLayerContext)
}

//----------------------------
// RenderLayerContext
//----------------------------

type RenderLayerContext interface {
	Path() PathPainter
	Rect() RectPainter
	Circle() CirclePainter
	Text() TextPainter

	Reset()
	Cleanup()
}

// Create a render context for the element with the specified 'elementId'.
// elementId should correspond to either an 'svg' or 'canvas' element.
func NewRenderContext(elementId string, layer RenderLayer) RenderContext {
	e := document.GetElementById(elementId)
	if e == nil {
		return nil
	}
	tag := strings.ToUpper(e.Tag)

	var context RenderContext
	switch tag {
	case "SVG", "G":
		context = NewSvgRenderContext(elementId)
	case "CANVAS":
		context = NewCanvasRenderContext(elementId)
	}
	if context == nil {
		return nil
	}
	if layer != nil {
		context.Layer(layer)
	}
	return context
}

//----------------------------
// FillLayer
// implements RenderLayer
//----------------------------

type FillLayer struct {
	Width, Height float64 //# 500,500
	Fill          string  // fill: #EEE
}

func NewFillLayer(width, height float64, fill string) *FillLayer {
	return &FillLayer{width, height, fill}
}

func (l *FillLayer) Render(context RenderLayerContext) {
	rectPainter := context.Rect()
	rectPainter.Rect(l.Width, l.Height)
	rectPainter.Fill(map[string]string{"fill": l.Fill})
}

//----------------------------
// Painter
//----------------------------

type Painter interface {
	Fill(style map[string]string)
}

type PathPainter interface {
	Painter

	Path(points []Vertex)
}

type RectPainter interface {
	Painter

	Rect(width, height float64)
}

type CirclePainter interface {
	Painter
}

type TextPainter interface {
	Painter

	FillText(m [6]float64, text string, style map[string]string)
}

//----------------------------
// SvgRenderContext
//----------------------------

type SvgRenderContext struct {
	layerAndContexts []svgRenderLayerAndContext
	svg              *document.Element
}

func NewSvgRenderContext(elementId string) RenderContext {
	c := &SvgRenderContext{}
	c.svg = document.GetElementById(elementId)
	return c
}

func (c *SvgRenderContext) Render() {
	c.Reset()
	for _, lc := range c.layerAndContexts {
		lc.context.Reset()
		lc.layer.Render(lc.context)
		lc.context.Cleanup()
	}
	c.Cleanup()
}

func (c *SvgRenderContext) Animate() RenderAnimator {
	return nil
}

func (c *SvgRenderContext) Layer(layer RenderLayer) {
	group := document.CreateElementNS(SVG_NS, "g")
	c.svg.AppendChild(group)
	lc := svgRenderLayerAndContext{layer, NewSvgLayerRenderContext(group)}
	c.layerAndContexts = append(c.layerAndContexts, lc)
}

func (c *SvgRenderContext) Reset() {

}

func (c *SvgRenderContext) Cleanup() {

}

//----------------------------
// svgRenderLayerAndContext
//----------------------------

type svgRenderLayerAndContext struct {
	layer   RenderLayer
	context RenderLayerContext
}

//----------------------------
// SvgLayerRenderContext
// implements RenderLayerContext
//----------------------------

type SvgLayerRenderContext struct {
	group         *document.Element
	pathPainter   PathPainter
	textPainter   TextPainter
	circlePainter CirclePainter
	rectPainter   RectPainter
	i             int
}

func NewSvgLayerRenderContext(group *document.Element) RenderLayerContext {
	c := &SvgLayerRenderContext{}
	c.group = group
	c.pathPainter = NewSvgPathPainter(c.elementFactory)
	c.textPainter = NewSvgTextPainter(c.elementFactory)
	c.circlePainter = NewSvgCirclePainter(c.elementFactory)
	c.rectPainter = NewSvgRectPainter(c.elementFactory)
	return c
}

// Returns an element with tagname `type`.
//
// This method uses an internal iterator to add new elements as they are
// drawn. If there is no child element at the current index, we append one
// and return it. If an element exists at the current index and it is the
// same tag, we return that. If the element is a different type, we create
// one and replace it then return it.
func (c *SvgLayerRenderContext) elementFactory(tag string) *document.Element {
	children := c.group.ChildNodes
	if c.i >= len(children) {
		path := document.CreateElementNS(SVG_NS, tag)
		c.group.AppendChild(path)
		c.i++
		return path
	}

	current := children[c.i]
	if current.Tag == tag {
		c.i++
		return current
	}

	path := document.CreateElementNS(SVG_NS, tag)
	c.group.ReplaceChild(path, current)
	c.i++
	return path
}

func (c *SvgLayerRenderContext) Path() PathPainter {
	return c.pathPainter
}

func (c *SvgLayerRenderContext) Rect() RectPainter {
	return c.rectPainter
}

func (c *SvgLayerRenderContext) Circle() CirclePainter {
	return c.circlePainter
}

func (c *SvgLayerRenderContext) Text() TextPainter {
	return c.textPainter
}

func (c *SvgLayerRenderContext) Reset() {
	c.i = 0
}

func (c *SvgLayerRenderContext) Cleanup() {
	children := c.group.ChildNodes
	for c.i < len(children) {
		children[c.i].SetAttribute("style", "display: none;")
		c.i++
	}
}

//----------------------------
// SvgPainter
//----------------------------

type SvgPainter struct {
	svgTag         string
	elementFactory func(tag string) *document.Element
	attributes     map[string]string
}

func NewSvgPainter(svgTag string, elementFactory func(tag string) *document.Element) *SvgPainter {
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

func (p *SvgPainter) Draw(style map[string]string) {
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

//----------------------------
// SvgPathPainter
//----------------------------

type SvgPathPainter struct {
	*SvgPainter
}

func NewSvgPathPainter(elementFactory func(tag string) *document.Element) *SvgPathPainter {
	return &SvgPathPainter{NewSvgPainter("path", elementFactory)}
}

func (p *SvgPathPainter) Path(points []Vertex) {
	str := "M"
	for _, point := range points {
		str += joinFloats(point.X, point.Y) + "L"
	}
	p.attributes["d"] = str[:len(str)-1]
}

//----------------------------
// SvgTextPainter
//----------------------------

type SvgTextPainter struct {
	*SvgPainter
}

func NewSvgTextPainter(elementFactory func(tag string) *document.Element) *SvgTextPainter {
	return &SvgTextPainter{&SvgPainter{svgTag: "text", elementFactory: elementFactory}}
}

func (p *SvgTextPainter) FillText(m [6]float64, text string, style map[string]string) {
	el := p.elementFactory(p.svgTag)

	// Mat3x4
	// | m0 m1  m2  m3 |
	// | m4 m5  m6  m7 |
	// | m8 m9 m10 m11 |

	// transform = matrix(<a> <b> <c> <d> <e> <f>)
	// | a  c  e |
	// | b  d  f |
	// | 0  0  1 |

	// set the transform attribute given the matrix m
	el.SetAttribute("transform", "matrix("+joinFloats(m[0], m[3], -m[1], -m[4], m[2], m[5])+")")

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

//----------------------------
// SvgCirclePainter
//----------------------------

type SvgCirclePainter struct {
	*SvgPainter
}

func NewSvgCirclePainter(elementFactory func(tag string) *document.Element) *SvgCirclePainter {
	return &SvgCirclePainter{NewSvgPainter("circle", elementFactory)}
}

//----------------------------
// SvgRectPainter
//----------------------------

type SvgRectPainter struct {
	*SvgPainter
}

func NewSvgRectPainter(elementFactory func(tag string) *document.Element) *SvgRectPainter {
	return &SvgRectPainter{NewSvgPainter("rect", elementFactory)}
}

func (p *SvgRectPainter) Rect(width, height float64) {
	p.attributes["width"] = formatFloat(width)
	p.attributes["height"] = formatFloat(height)
}

////////
func NewCanvasRenderContext(elementId string) RenderContext {
	return nil
}
