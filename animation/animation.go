package animation

import (
	"time"

	"github.com/vibrantgio/seen"
)

type Handler func(t, dt time.Duration)

type Animator interface {
	Start() Animator
	Stop()
	OnBefore(handler Handler) Animator
	OnFrame(handler Handler) Animator
	OnAfter(handler Handler) Animator
}

// Animation is an animation loop running at 33Hz.
// We supply pre and post events for applying animation
// changes between frames.
type Animation struct {
	ticker   seen.IntervalID
	onBefore []Handler
	onFrame  []Handler
	onAfter  []Handler
}

// Start makes an animation loop call the OnFrame handlers at 33 Hz.
func (a *Animation) Start() Animator {
	if a.ticker != 0 {
		seen.ClearInterval(a.ticker)
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
	a.ticker = seen.SetInterval(animate, 30*time.Millisecond)
	return a
}

// Stop makes the animation stop calling the OnFrame handlers.
func (a *Animation) Stop() {
	seen.ClearInterval(a.ticker)
}

// OnBefore handlers are called before the OnFrame handlers are called.
func (a *Animation) OnBefore(handler Handler) Animator {
	a.onBefore = append(a.onBefore, handler)
	return a
}

// OnFrame handlers are called at a frequency of 33Hz.
func (a *Animation) OnFrame(handler Handler) Animator {
	a.onFrame = append(a.onFrame, handler)
	return a
}

// OnAfter handlers are called after every OnFrame handler call.
func (a *Animation) OnAfter(handler Handler) Animator {
	a.onAfter = append(a.onAfter, handler)
	return a
}
