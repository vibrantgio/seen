package drag

import (
	"math"
	"time"

	"github.com/vibrantgio/seen"
)

type Type string

const (
	Start      = Type("Start")
	Move       = Type("Move")
	End        = Type("End")
	EndInertia = Type("EndInertia")
)

type Event struct {
	Type Type
	X    float64
	Y    float64
	T    time.Duration
	Dx   float64
	Dy   float64
	Dt   time.Duration
}

type Handler func(Event)

type Dragger interface {
	On(handler Handler)
}

type Option func(*Drag)

func Inertia(value bool) Option {
	return func(d *Drag) {
		if value {
			motion := DefaultInertialMotion
			d.inertia = &motion
		} else {
			d.inertia = nil
		}
	}
}

type Drag struct {
	handlers []Handler
	inertia  *InertialMotion
	iid      seen.IntervalID
}

func DragWith(options ...Option) *Drag {
	drag := &Drag{}
	for _, option := range options {
		option(drag)
	}
	return drag
}

func (d *Drag) On(handler Handler) {
	d.handlers = append(d.handlers, handler)
}

func (d *Drag) Handle(e Event) {
	if d.inertia != nil {
		// Handle dragging with inertia
		InertiaMove := func(t, dt time.Duration) bool {
			dx, dy := d.inertia.Damp().Get()
			if math.Abs(dx) < 1 && math.Abs(dy) < 1 {
				e.Type = EndInertia
				e.Dx, e.Dy, e.Dt = 0, 0, 0
				d.Notify(e)
				return false
			}
			e.Type = Move
			e.X += dx
			e.Y += dy
			e.T += dt
			e.Dx, e.Dy, e.Dt = dx, dy, dt
			d.Notify(e)
			return true
		}
		switch e.Type {
		case Start:
			d.inertia.Reset()
			seen.ClearInterval(d.iid)
		case Move:
			d.inertia.Update(e.Dx, e.Dy, e.Dt)
		case End:
			d.iid = seen.SetInterval(InertiaMove, d.inertia.InertiaDelay)
		}
	}
	d.Notify(e)
}

func (d *Drag) Notify(e Event) {
	for _, handler := range d.handlers {
		handler(e)
	}
}
