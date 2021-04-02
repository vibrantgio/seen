package canvas

import (
	"strings"

	"github.com/reactivego/seen"
	"github.com/reactivego/seen/document"
	"github.com/reactivego/seen/render"
	"github.com/reactivego/seen/render/svg"
)

type Context struct{ canvas *document.Element }

// MakeContext creates a render context for the element with the specified 'elementId'.
// The elementId should correspond to a 'canvas' element.
func MakeContext(elementId string, layer render.RenderLayer) render.RenderContext {
	e := document.GetElementById(elementId)
	if e == nil {
		return nil
	}
	tag := strings.ToUpper(e.Tag)
	if tag == "SVG" || tag == "G" {
		return svg.MakeContext(elementId, layer)
	}
	var context render.RenderContext
	if tag == "CANVAS" {
		context = &Context{canvas: document.GetElementById(elementId)}
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
}

func (c *Context) Render() {
}

func (c *Context) Animate() seen.Animator {
	animator := seen.MakeAnimator()
	animator.OnFrame(func(d, dt float64) { c.Render() })
	return animator
}
