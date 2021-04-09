package seen

import "time"

// Animator is a class that is useful for creating an animation loop.
// We supply pre and post events for applying animation changes between frames.
type Animator struct {
	ticker   IntervalID
	onBefore []Handler
	onFrame  []Handler
	onAfter  []Handler
}

type Handler func(t, dt time.Duration)

// Start makes the animator call the OnFrame handlers at 33 Hz.
func (a *Animator) Start() {
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

// Stop makes the animator stop calling the OnFrame handlers.
func (a *Animator) Stop() {
	ClearInterval(a.ticker)
}

// OnBefore handlers are called before the OnFrame handlers are called.
func (a *Animator) OnBefore(handler Handler) {
	a.onBefore = append(a.onBefore, handler)
}

// OnFrame handlers are called at a frequency of 33Hz.
func (a *Animator) OnFrame(handler Handler) {
	a.onFrame = append(a.onFrame, handler)
}

// OnAfter handlers are called after every OnFrame handler call.
func (a *Animator) OnAfter(handler Handler) {
	a.onAfter = append(a.onAfter, handler)
}
