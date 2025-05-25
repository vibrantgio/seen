package svg

import (
	"strings"

	"github.com/vibrantgio/seen/animation"
	"github.com/vibrantgio/seen/context"
	"github.com/vibrantgio/seen/drag"
	"github.com/vibrantgio/seen/layer"
	"github.com/vibrantgio/seen/zoom"
)

// Context
type Context struct {
	svg    *Element
	render []func()
}

var _ context.Context = (*Context)(nil)

// NewContext creates a render context to the element with the
// specified 'elementId'. This element should be an 'svg' element.
func NewContext(element *Element, layers ...layer.Layer) *Context {
	if element == nil {
		return nil
	}
	var context *Context
	tag := strings.ToUpper(element.Tag)
	if tag != "SVG" && tag != "G" {
		return nil
	}
	context = &Context{svg: element}
	context.SetLayers(layers...)
	return context
}

func (c *Context) SetLayers(layers ...layer.Layer) {
	c.render = nil
	for _, layer := range layers {
		layer := layer // no longer needed when go 1.22 is set in go.mod
		group := c.svg.CreateElementNS(SVG_NS, "g")
		c.svg.AppendChild(group)
		canvas := NewCanvas(group)
		c.render = append(c.render, func() {
			canvas.Reset()
			layer.RenderOn(canvas)
			canvas.Cleanup()
		})
	}
}

func (c *Context) Render() {
	for _, render := range c.render {
		render()
	}
}

func (c *Context) Animate() animation.Animator {
	return nil
}

func (c *Context) Drag(...drag.Option) drag.Dragger {
	return nil
}

func (c *Context) Zoom(...zoom.Option) zoom.Zoomer {
	return nil
}
