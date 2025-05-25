package zoom

import "math"

type Type string

const Move = Type("Zoom")

type Event struct {
	Type       Type
	Dx         float64
	Dy         float64
	WheelDelta float64
	Zoom       float64
}

type Handler func(Event)

type Zoomer interface {
	On(handler Handler)
}

type Option func() float64

func Speed(speed float64) Option {
	return func() float64 { return speed }
}

type Zoom struct {
	Speed float64

	handlers []Handler
}

func With(options ...Option) *Zoom {
	zoom := &Zoom{Speed: 0.25}
	for _, option := range options {
		zoom.Speed = option()
	}
	return zoom
}

func (z *Zoom) On(handler Handler) {
	z.handlers = append(z.handlers, handler)
}

func (z *Zoom) Handle(e Event) {
	e.Zoom = math.Exp2(e.WheelDelta / 120 * z.Speed)
	for _, handler := range z.handlers {
		handler(e)
	}
}
