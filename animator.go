package seen

// The animator class is useful for creating an animation loop. We supply pre
// and post events for applying animation changes between frames.
type Animator interface {
	// Start makes the animator call the Render method at 33 Hz.
	Start()
	// Stop makes the animator stop calling the Render method.
	Stop()

	// OnBefore is called before every Render call. Hook this to change the
	// scene.
	OnBefore(Handler)

	// OnFrame is called at 33Hz.
	OnFrame(Handler)

	// OnAfter is called after every Render call.
	OnAfter(func())
}

type Handler func(t, dt float64)

func MakeAnimator() Animator {
	return &animator{}
}

type animator struct {
	onBefore Handler
	handler  Handler
	onAfter  func()
}

func (a *animator) Start() {

}

func (a *animator) Stop() {

}

func (a *animator) OnBefore(cb Handler) {
	a.onBefore = cb
}

func (a *animator) OnFrame(handler Handler) {
	a.handler = handler
}

func (a *animator) OnAfter(cb func()) {
	a.onAfter = cb
}
