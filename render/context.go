package render

import (
	"strings"
	"github.com/reactivego/seen/document"
	"github.com/reactivego/seen"
)

// RenderContext
type RenderContext interface {
	Render()
	Animate() seen.RenderAnimator
	Layer(layer RenderLayer)

	Reset()
	Cleanup()
}

// MakeRenderContext creates a render context for the element with the specified 'elementId'.
// elementId should correspond to either an 'svg' or 'canvas' element.
func MakeRenderContext(elementId string, layer RenderLayer) RenderContext {
	e := document.GetElementById(elementId)
	if e == nil {
		return nil
	}
	tag := strings.ToUpper(e.Tag)
	var context RenderContext
	switch tag {
	case "SVG", "G":
		context = MakeSvgRenderContext(elementId)
	case "CANVAS":
		context = MakeCanvasRenderContext(elementId)
	}
	if context == nil {
		return nil
	}
	if layer != nil {
		context.Layer(layer)
	}
	return context
}
