package render

import "github.com/reactivego/seen"

// RenderAnimator controls whether the Render method on the RenderContext is
// called at 33 Hz.
func MakeRenderAnimator(context RenderContext) seen.Animator {
	animator := seen.MakeAnimator()
	animator.OnFrame(func(d, dt float64) {
		context.Render()
	})
	return animator
}
