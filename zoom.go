package seen

import "math"

type ZoomType string

const ZoomMove = ZoomType("Zoom")

type ZoomEvent struct {
	Type       ZoomType
	Dx         float64
	Dy         float64
	WheelDelta float64
	Zoom       float64
}

type ZoomHandler func(ZoomEvent)

type Zoomer interface {
	On(handler ZoomHandler)
}

type ZoomOption func() float64

func Speed(speed float64) ZoomOption {
	return func() float64 { return speed }
}

type Zoom struct {
	Speed float64

	handlers []ZoomHandler
}

func ZoomWith(options ...ZoomOption) *Zoom {
	zoom := &Zoom{Speed: 0.25}
	for _, option := range options {
		zoom.Speed = option()
	}
	return zoom
}

func (z *Zoom) On(handler ZoomHandler) {
	z.handlers = append(z.handlers, handler)
}

func (z *Zoom) Handle(e ZoomEvent) {
	e.Zoom = math.Exp2(e.WheelDelta / 120 * z.Speed)
	for _, handler := range z.handlers {
		handler(e)
	}
}
