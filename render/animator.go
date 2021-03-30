package render

// RenderAnimator controls whether the Render method on the RenderContext is
// called at 33 Hz.
type RenderAnimator interface {
	// Start makes the animator call the Render method at 33 Hz.
	Start()
	// Stop makes the animator stop calling the Render method.
	Stop()

	// OnBefore is called before every Render call. Hook this to change the
	// scene.
	OnBefore(func(d,dt float64))

	// OnAfter is called after every Render call.
	OnAfter(func())
}

type Animator struct {
	onBefore func(d,dt float64)
	onAfter func()
}

func MakeAnimator() RenderAnimator {
	return &Animator{}
}

func (a *Animator) Start() {

}

func (a *Animator) Stop() {

}

func (a * Animator) OnBefore(cb func(t,dt float64)) {
	a.onBefore = cb
}

func (a * Animator) OnAfter(cb func()) {
	a.onAfter = cb
}
