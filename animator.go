package seen

type RenderAnimator interface {
	OnFrame(handler func (t, dt float64))
	Start()
}

type Animator struct {
	onFrameHandler func (t, dt float64)
}

func MakeAnimator() RenderAnimator {
	return &Animator{}
}

func (a *Animator) OnFrame(handler func (t, dt float64)) {
	a.onFrameHandler = handler
}

func (a *Animator) Start() {

}
