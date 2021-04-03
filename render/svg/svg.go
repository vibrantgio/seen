package svg

import (
	"strings"

	"github.com/reactivego/seen"
	"github.com/reactivego/seen/affine"
	"github.com/reactivego/seen/document"
	"github.com/reactivego/seen/render"
)

// Context
type Context struct {
	svg    *document.Element
	render []func()
}

// MakeContext creates a render context for the element with the
// specified 'elementId'. This element should be an 'svg' element.
func MakeContext(elementId string, layer render.RenderLayer) render.RenderContext {
	e := document.GetElementById(elementId)
	if e == nil {
		return nil
	}
	tag := strings.ToUpper(e.Tag)
	var context render.RenderContext
	if tag == "SVG" || tag == "G" {
		context = &Context{svg: document.GetElementById(elementId)}
	}
	if context == nil {
		return nil
	}
	if layer != nil {
		context.Layer(layer)
	}
	return context
}

func (c *Context) Layer(layer render.RenderLayer) {
	group := document.CreateElementNS(document.SVG_NS, "g")
	c.svg.AppendChild(group)
	painter := MakeSvgPainter(group)
	c.render = append(c.render, func() {
		painter.Reset()
		layer.Paint(painter)
		painter.Cleanup()
	})
}

func (c *Context) Render() {
	for _, render := range c.render {
		render()
	}
}

func (c *Context) Animate() seen.Animator {
	return nil
}

// SvgPainter
type SvgPainter struct {
	group         *document.Element
	pathPainter   render.PathPainter
	textPainter   render.TextPainter
	circlePainter render.CirclePainter
	rectPainter   render.RectPainter
	i             int
}

func MakeSvgPainter(group *document.Element) render.Painter {
	c := &SvgPainter{}
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
func (c *SvgPainter) elementFactory(tag string) *document.Element {
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

func (c *SvgPainter) Path() render.PathPainter {
	return c.pathPainter
}

func (c *SvgPainter) Rect() render.RectPainter {
	return c.rectPainter
}

func (c *SvgPainter) Circle() render.CirclePainter {
	return c.circlePainter
}

func (c *SvgPainter) Text() render.TextPainter {
	return c.textPainter
}

func (c *SvgPainter) Reset() {
	c.i = 0
}

func (c *SvgPainter) Cleanup() {
	children := c.group.ChildNodes
	for c.i < len(children) {
		children[c.i].SetAttribute("style", "display: none;")
		c.i++
	}
}

// SvgCommonPainter
type SvgCommonPainter struct {
	svgTag         string
	elementFactory func(tag string) *document.Element
	attributes     map[string]string
}

func MakeSvgCommonPainter(svgTag string, elementFactory func(tag string) *document.Element) *SvgCommonPainter {
	p := &SvgCommonPainter{}
	if svgTag != "" {
		p.svgTag = svgTag
	} else {
		p.svgTag = "g"
	}
	p.elementFactory = elementFactory
	p.attributes = make(map[string]string)
	return p
}

func (p *SvgCommonPainter) Clear() {
	p.attributes = make(map[string]string)
}

func (p *SvgCommonPainter) Fill(style map[string]string) {
	p.Paint(style)
}

func (p *SvgCommonPainter) Stroke(style map[string]string) {
	p.Paint(style)
}

func (p *SvgCommonPainter) Paint(style map[string]string) {
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
	*SvgCommonPainter
}

func MakeSvgPathPainter(elementFactory func(tag string) *document.Element) *SvgPathPainter {
	return &SvgPathPainter{MakeSvgCommonPainter("path", elementFactory)}
}

func (p *SvgPathPainter) Path(points []seen.Point) {
	str := "M"
	for _, point := range points {
		str += render.Fjoin(point.X, point.Y) + "L"
	}
	p.attributes["d"] = str[:len(str)-1]
}

// SvgTextPainter
type SvgTextPainter struct {
	*SvgCommonPainter
}

func MakeSvgTextPainter(elementFactory func(tag string) *document.Element) *SvgTextPainter {
	return &SvgTextPainter{&SvgCommonPainter{svgTag: "text", elementFactory: elementFactory}}
}

func (p *SvgTextPainter) FillText(t affine.Matrix, text string, style render.Style) {
	el := p.elementFactory(p.svgTag)

	// set the transform attribute given the matrix m
	el.SetAttribute("transform", "matrix("+render.Fjoin(t.A, t.B, t.C, t.D, t.E, t.F)+")")

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
	*SvgCommonPainter
}

func MakeSvgCirclePainter(elementFactory func(tag string) *document.Element) *SvgCirclePainter {
	return &SvgCirclePainter{MakeSvgCommonPainter("circle", elementFactory)}
}

// SvgRectPainter
type SvgRectPainter struct {
	*SvgCommonPainter
}

func MakeSvgRectPainter(elementFactory func(tag string) *document.Element) *SvgRectPainter {
	return &SvgRectPainter{MakeSvgCommonPainter("rect", elementFactory)}
}

func (p *SvgRectPainter) Size(width, height float64) {
	p.attributes["width"] = render.Ftoa(width)
	p.attributes["height"] = render.Ftoa(height)
}

func (p *SvgRectPainter) CornerRadius(rx, ry float64) {
	p.attributes["rx"] = render.Ftoa(rx)
	p.attributes["ry"] = render.Ftoa(ry)
}
