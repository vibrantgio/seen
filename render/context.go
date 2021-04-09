package render

import "github.com/reactivego/seen"

// RenderContext
type RenderContext interface {
	Layer(RenderLayer)
	Render()
	Animate() *seen.Animator
}
