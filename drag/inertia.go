package drag

import (
	"math"
	"time"
)

type InertialMotion struct {
	// InertiaExtinction is the amount taken off the remaining internal motion on every
	// call to the Damp() method. A value of e.g. 0.1 means that 10% is taken off.
	InertiaExtinction float64

	// SmoothingTimeout is the duration of the smoothing low pass filter.
	SmoothingTimeout time.Duration

	// InertiaDelay is the expected time between calls to Get to retrieve the
	// remaining inertial motion.
	InertiaDelay time.Duration

	// dx,dy contains the inertial 2D motion in pixels per millisecond.
	dx, dy float64
}

var DefaultInertialMotion = InertialMotion{
	InertiaExtinction: 0.1,
	SmoothingTimeout:  300 * time.Millisecond,
	InertiaDelay:      30 * time.Millisecond,
}

// Get returns the residual inertial motion that occured in a period of duration
// InertiaMsecDelay. So it is expected to be called every InertiaMsecDelay
// milliseconds.
func (i *InertialMotion) Get() (dx, dy float64) {
	scale := 1000.0 / float64(i.InertiaDelay.Milliseconds())
	return i.dx * scale, i.dy * scale
}

// Reset zeroes the residual inertial motion.
func (i *InertialMotion) Reset() {
	i.dx, i.dy = 0, 0
}

// Update will update the current rate of motion dx,dy in x and y.
// The dt value is expected to be in milliseconds and it specifies the duration
// of the period in which dx,dy were accrued.
func (i *InertialMotion) Update(dx, dy float64, dt time.Duration) {
	dtms := float64(dt.Milliseconds())
	stms := float64(i.SmoothingTimeout.Milliseconds())
	// Pixels per milliseconds
	dx, dy = dx/math.Max(dtms, 1), dy/math.Max(dtms, 1)
	// On initial update take dx,dy as is.
	if i.dx == 0 && i.dy == 0 {
		dtms = stms
	}
	// Smoothing interpolation based on time between measurements
	t := math.Min(1, dtms/stms)
	i.dx, i.dy = t*dx+(1.0-t)*i.dx, t*dy+(1.0-t)*i.dy
}

// Apply damping to slow the motion once the user has stopped dragging.
func (i *InertialMotion) Damp() *InertialMotion {
	i.dx, i.dy = i.dx*(1-i.InertiaExtinction), i.dy*(1-i.InertiaExtinction)
	return i
}
