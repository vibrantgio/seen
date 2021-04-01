package canvas

import (
	"strings"

	"github.com/reactivego/seen"
	"github.com/reactivego/seen/affine"
	"github.com/reactivego/seen/document"
	"github.com/reactivego/seen/render"
	"github.com/reactivego/seen/render/svg"
)

// MakeRenderContext creates a render context for the element with the specified 'elementId'.
// The elementId should correspond to a 'canvas' element.
func MakeRenderContext(elementId string, layer render.RenderLayer) render.RenderContext {
	e := document.GetElementById(elementId)
	if e == nil {
		return nil
	}
	tag := strings.ToUpper(e.Tag)
	if tag == "SVG" || tag == "G" {
		return svg.MakeRenderContext(elementId, layer)
	}
	var context render.RenderContext
	if tag == "CANVAS" {
		context = &CanvasRenderContext{canvas: document.GetElementById(elementId)}
	}
	if context == nil {
		return nil
	}
	if layer != nil {
		context.Layer(layer)
	}
	return context
}

//TODO

type CanvasRenderContext struct {
	canvas *document.Element
}

func (c *CanvasRenderContext) Render() {

}

func (c *CanvasRenderContext) Animate() seen.Animator {
	return render.MakeRenderAnimator(c)
}

func (c *CanvasRenderContext) Layer(layer render.RenderLayer) {

}

func (c *CanvasRenderContext) Reset() {

}

func (c *CanvasRenderContext) Cleanup() {

}

// CanvasPainter
type CanvasPainter struct {
}

// CanvasTextPainter
type CanvasTextPainter struct {
	*CanvasPainter
}

func MakeCanvasTextPainter() *CanvasTextPainter {
	return &CanvasTextPainter{&CanvasPainter{}}
}

func (p *CanvasTextPainter) FillText(t *affine.Matrix, text string, style render.Style) {
}
