package render

import "github.com/reactivego/seen"

// MakeRenderAnimator returns an Animator that calls the context.Render method 33 times per second.
func MakeRenderAnimator(context RenderContext) seen.Animator {
	animator := seen.MakeAnimator()
	animator.OnFrame(func(d, dt float64) {
		context.Render()
	})
	return animator
}
