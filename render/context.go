package render

import "github.com/reactivego/seen"

// RenderContext
type RenderContext interface {
	Layers(...RenderLayer)
	Render()
	Animate() *seen.Animator
	Drag() *seen.Drag
	Zoom() *seen.Zoom
}
