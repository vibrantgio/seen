package render

import "github.com/reactivego/seen"

// RenderContext
type RenderContext interface {
	Render()
	Animate() seen.Animator
	Layer(RenderLayer)

	Reset()
	Cleanup()
}
