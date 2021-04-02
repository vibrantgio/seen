package seen

type Drag struct {
	Inertia bool

	handlers []DragHandler
}

type DragHandler func(DragEvent)

type DragEvent struct {
	OffsetRelativeX, OffsetRelativeY float64
}

type DragOption func(*bool)

func Inertia(value bool) DragOption {
	return func(v *bool) {
		*v = value
	}
}

func MakeDrag(options ...DragOption) *Drag {
	drag := &Drag{}
	for _, option := range options {
		option(&drag.Inertia)
	}
	return drag
}

func (d *Drag) On(handler DragHandler) {
	d.handlers = append(d.handlers, handler)
}

func (d *Drag) Handle(e DragEvent) {
	for _, handler := range d.handlers {
		handler(e)
	}
}
