package seen

import (
	"math"
	"time"
)

type DragType string

const (
	DragStart      = DragType("DragStart")
	DragMove       = DragType("Drag")
	DragEnd        = DragType("DragEnd")
	DragEndInertia = DragType("DragEndInertia")
)

type DragEvent struct {
	Type DragType
	X    float64
	Y    float64
	T    time.Duration
	Dx   float64
	Dy   float64
	Dt   time.Duration
}

type DragHandler func(DragEvent)

type Dragger interface {
	On(handler DragHandler)
}

type DragOption func(*Drag)

func Inertia(value bool) DragOption {
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
	handlers []DragHandler
	inertia  *InertialMotion
	iid      IntervalID
}

func DragWith(options ...DragOption) *Drag {
	drag := &Drag{}
	for _, option := range options {
		option(drag)
	}
	return drag
}

func (d *Drag) On(handler DragHandler) {
	d.handlers = append(d.handlers, handler)
}

func (d *Drag) Handle(e DragEvent) {
	if d.inertia != nil {
		// Handle dragging with inertia
		InertiaMove := func(t, dt time.Duration) bool {
			dx, dy := d.inertia.Damp().Get()
			if math.Abs(dx) < 1 && math.Abs(dy) < 1 {
				e.Type = DragEndInertia
				e.Dx, e.Dy, e.Dt = 0, 0, 0
				d.Notify(e)
				return false
			}
			e.Type = DragMove
			e.X += dx
			e.Y += dy
			e.T += dt
			e.Dx, e.Dy, e.Dt = dx, dy, dt
			d.Notify(e)
			return true
		}
		switch e.Type {
		case DragStart:
			d.inertia.Reset()
			ClearInterval(d.iid)
		case DragMove:
			d.inertia.Update(e.Dx, e.Dy, e.Dt)
		case DragEnd:
			d.iid = SetInterval(InertiaMove, d.inertia.InertiaDelay)
		}
	}
	d.Notify(e)
}

func (d *Drag) Notify(e DragEvent) {
	for _, handler := range d.handlers {
		handler(e)
	}
}
