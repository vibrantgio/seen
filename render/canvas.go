package render

import (
	"github.com/reactivego/seen/document"
	"github.com/reactivego/seen"
	"github.com/reactivego/seen/affine"
)

func MakeCanvasRenderContext(elementId string) RenderContext {
	c := &CanvasRenderContext{}
	c.canvas = document.GetElementById(elementId)
	return c
}

//TODO

type CanvasRenderContext struct {
	canvas *document.Element
}

func (c *CanvasRenderContext) Render() {

}

func (c *CanvasRenderContext) Animate() seen.RenderAnimator {
	return seen.MakeAnimator()
}

func (c *CanvasRenderContext) Layer(layer RenderLayer) {

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

func (p *CanvasTextPainter) FillText(t *affine.Matrix, text string, style map[string]string) {
	// p.ctx.save()
	// p.ctx.setTransform(m[0], m[3], -m[1], -m[4], m[2], m[5])

	// if style.font? then @ctx.font = style.font
	// if style.fill? then @ctx.fillStyle = style.fill
	// if style['text-anchor']? then @ctx.textAlign = @_cssToCanvasAnchor(style['text-anchor'])

	// p.ctx.fillText(text, 0, 0)
	// p.ctx.restore()
	// return @
}
