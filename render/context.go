package render

import "github.com/reactivego/seen"

// RenderContext
type RenderContext interface {
	Layers(...Layer)
	Render()
	Animate() *seen.Animator
	Drag() *seen.Drag
	Zoom() *seen.Zoom
}
