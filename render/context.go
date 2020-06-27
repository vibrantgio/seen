package render

// RenderContext
type RenderContext interface {
	Render()
	Animate() RenderAnimator
	Layer(RenderLayer)

	Reset()
	Cleanup()
}
