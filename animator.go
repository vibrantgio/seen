package seen

import "time"

// The animator class is useful for creating an animation loop. We supply pre
// and post events for applying animation changes between frames.
type Animator interface {
	// Start makes the animator call the OnFrame handlers at 33 Hz.
	Start()

	// Stop makes the animator stop calling the OnFrame handlers.
	Stop()

	// OnBefore handlers are called before the OnFrame handlers are called.
	OnBefore(Handler)

	// OnFrame handlers are called at a frequency of 33Hz.
	OnFrame(Handler)

	// OnAfter handlers are called after every OnFrame handler call.
	OnAfter(Handler)
}

type Handler func(t, dt time.Duration)

type animator struct {
	ticker   IntervalID
	onBefore []Handler
	onFrame  []Handler
	onAfter  []Handler
}

func MakeAnimator() Animator {
	return &animator{}
}

func (a *animator) Start() {
	if a.ticker != 0 {
		ClearInterval(a.ticker)
	}
	animate := func(t, dt time.Duration) bool {
		for _, handler := range a.onBefore {
			handler(t, dt)
		}
		for _, handler := range a.onFrame {
			handler(t, dt)
		}
		for _, handler := range a.onAfter {
			handler(t, dt)
		}
		return true
	}
	a.ticker = SetInterval(animate, 30*time.Millisecond)
}

func (a *animator) Stop() {
	ClearInterval(a.ticker)
}

func (a *animator) OnBefore(handler Handler) {
	a.onBefore = append(a.onBefore, handler)
}

func (a *animator) OnFrame(handler Handler) {
	a.onFrame = append(a.onFrame, handler)
}

func (a *animator) OnAfter(handler Handler) {
	a.onAfter = append(a.onAfter, handler)
}
