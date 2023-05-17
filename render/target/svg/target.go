package svg

import (
	"strings"

	"github.com/reactivego/seen"
	"github.com/reactivego/seen/document"
	"github.com/reactivego/seen/render"
)

// Target
type Target struct {
	svg    *document.Element
	render []func()
}

// NewTarget creates a render target to the element with the
// specified 'elementId'. This element should be an 'svg' element.
func NewTarget(element *document.Element, layers ...render.Layer) *Target {
	if element == nil {
		return nil
	}
	var context *Target
	tag := strings.ToUpper(element.Tag)
	if tag != "SVG" && tag != "G" {
		return nil
	}
	context = &Target{svg: element}
	context.SetLayers(layers...)
	return context
}

func (c *Target) SetLayers(layers ...render.Layer) {
	c.render = nil
	add := func(layer render.Layer) {
		group := c.svg.CreateElementNS(document.SVG_NS, "g")
		c.svg.AppendChild(group)
		painter := NewPainter(group)
		c.render = append(c.render, func() {
			painter.Reset()
			layer.Paint(painter)
			painter.Cleanup()
		})
	}
	for _, layer := range layers {
		add(layer)
	}
}

func (c *Target) Render() {
	for _, render := range c.render {
		render()
	}
}

func (c *Target) Animate() *seen.Animation {
	return nil
}

func (c *Target) Drag() *seen.Drag {
	return nil
}

func (c *Target) Zoom() *seen.Zoom {
	return nil
}
