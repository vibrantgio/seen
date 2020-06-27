package seen

// The animator class is useful for creating an animation loop. We supply pre
// and post events for apply animation changes between frames.
type Animator interface {
	OnFrame(Handler)

	Start()
	Stop()
}

type Handler func(t, dt float64)

func MakeAnimator() Animator {
	return &animator{}
}

type animator struct {
	handler Handler
}

func (a *animator) OnFrame(handler Handler) {
	a.handler = handler
}

func (a *animator) Start() {

}

func (a *animator) Stop() {
	
}
