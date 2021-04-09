package seen

import "math"

type Zoom struct {
	Speed float64

	handlers []ZoomHandler
}

type ZoomHandler func(ZoomEvent)

type ZoomEvent struct {
	Type       ZoomType
	Dx         float64
	Dy         float64
	WheelDelta float64
	Zoom       float64
}

type ZoomType string

const ZoomMove = ZoomType("Zoom")

type ZoomOption func() float64

func Speed(speed float64) ZoomOption {
	return func() float64 { return speed }
}

func ZoomWith(options ...ZoomOption) *Zoom {
	speed := 0.25
	for _, option := range options {
		speed = option()
	}
	return &Zoom{Speed: speed}
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
