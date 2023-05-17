package seen

import "time"

type AnimationHandler func(t, dt time.Duration)

type Animator interface {
	Start()
	Stop()
	OnBefore(handler AnimationHandler)
	OnFrame(handler AnimationHandler)
	OnAfter(handler AnimationHandler)
}

// Animation is an animation loop running at 33Hz.
// We supply pre and post events for applying animation
// changes between frames.
type Animation struct {
	ticker   IntervalID
	onBefore []AnimationHandler
	onFrame  []AnimationHandler
	onAfter  []AnimationHandler
}

// Start makes an animation loop call the OnFrame handlers at 33 Hz.
func (a *Animation) Start() {
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

// Stop makes the animation stop calling the OnFrame handlers.
func (a *Animation) Stop() {
	ClearInterval(a.ticker)
}

// OnBefore handlers are called before the OnFrame handlers are called.
func (a *Animation) OnBefore(handler AnimationHandler) {
	a.onBefore = append(a.onBefore, handler)
}

// OnFrame handlers are called at a frequency of 33Hz.
func (a *Animation) OnFrame(handler AnimationHandler) {
	a.onFrame = append(a.onFrame, handler)
}

// OnAfter handlers are called after every OnFrame handler call.
func (a *Animation) OnAfter(handler AnimationHandler) {
	a.onAfter = append(a.onAfter, handler)
}
