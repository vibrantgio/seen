package gio

import (
	"image"
	"math"
	"time"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/op"

	"github.com/vibrantgio/seen"
	"github.com/vibrantgio/seen/animation"
	"github.com/vibrantgio/seen/context"
	"github.com/vibrantgio/seen/context/gio/canvas"
	"github.com/vibrantgio/seen/drag"
	"github.com/vibrantgio/seen/layer"
	"github.com/vibrantgio/seen/zoom"
)

// Context
type Context struct {
	window   *app.Window
	layers   []layer.Layer
	inputs   []func(*op.Ops)
	handlers []func(event.Queue)
}

var _ context.Context = (*Context)(nil)

// NewContext creates a render context for the given gio app.Window.
func NewContext(window *app.Window, layers ...layer.Layer) *Context {
	t := &Context{window: window}
	t.SetLayers(layers...)
	return t
}

func (c *Context) SetLayers(layers ...layer.Layer) {
	c.layers = layers
}

func (c *Context) Render() {
	// Calling this will result in Process being called asynchronically
	c.window.Invalidate()
}

func (c *Context) Process(ops *op.Ops, queue event.Queue) {
	canvas := &canvas.Canvas{Ops: ops}
	for _, layer := range c.layers {
		layer.RenderOn(canvas)
	}
	for _, input := range c.inputs {
		input(ops)
	}
	for _, handler := range c.handlers {
		handler(queue)
	}
	if seen.Scheduler.Run() {
		op.InvalidateOp{}.Add(ops)
	}
}

func (c *Context) Animate() animation.Animator {
	animation := &animation.Animation{}
	animation.OnFrame(func(d, dt time.Duration) {
		c.Render()
	})
	return animation
}

func (c *Context) Drag(options ...drag.Option) drag.Dragger {
	d := drag.DragWith(options...)
	c.inputs = append(c.inputs, func(ops *op.Ops) {
		defer pointer.PassOp{}.Push(ops).Pop()
		const types = pointer.Press | pointer.Drag | pointer.Release
		pointer.InputOp{Tag: d, Types: types}.Add(ops)
	})
	previous := struct {
		Position f32.Point
		Time     time.Duration
	}{}
	c.handlers = append(c.handlers, func(q event.Queue) {
		for _, event := range q.Events(d) {
			if p, ok := event.(pointer.Event); ok {
				switch p.Type {
				case pointer.Press:
					d.Handle(drag.Event{
						Type: drag.Start,
						X:    float64(p.Position.X),
						Y:    float64(p.Position.Y),
						T:    p.Time,
					})
				case pointer.Drag:
					if previous.Time != 0 {
						dP := p.Position.Sub(previous.Position)
						dT := p.Time - previous.Time
						d.Handle(drag.Event{
							Type: drag.Move,
							X:    float64(p.Position.X),
							Y:    float64(p.Position.Y),
							T:    p.Time,
							Dx:   float64(dP.X),
							Dy:   float64(dP.Y),
							Dt:   dT,
						})
					}
					previous.Position, previous.Time = p.Position, p.Time
				case pointer.Release:
					if previous.Time != 0 {
						dP := p.Position.Sub(previous.Position)
						dT := p.Time - previous.Time
						d.Handle(drag.Event{
							Type: drag.End,
							X:    float64(p.Position.X),
							Y:    float64(p.Position.Y),
							T:    p.Time,
							Dx:   float64(dP.X),
							Dy:   float64(dP.Y),
							Dt:   dT,
						})
					}
					previous.Time = 0
				}
			}
		}
	})
	return d
}

func (c *Context) Zoom(options ...zoom.Option) zoom.Zoomer {
	z := zoom.With(options...)
	c.inputs = append(c.inputs, func(ops *op.Ops) {
		pointer.PassOp{}.Push(ops).Pop()
		pointer.InputOp{
			Tag:          z,
			Types:        pointer.Scroll,
			ScrollBounds: image.Rect(-120, -120, 120, 120),
		}.Add(ops)
	})
	c.handlers = append(c.handlers, func(q event.Queue) {
		for _, event := range q.Events(z) {
			if p, ok := event.(pointer.Event); ok {
				dx, dy := -float64(p.Scroll.X), -float64(p.Scroll.Y)
				dxy := math.Copysign(math.Hypot(dx, dy), dy)
				z.Handle(zoom.Event{
					Type:       zoom.Move,
					Dx:         dx,
					Dy:         dy,
					WheelDelta: dxy,
				})
			}
		}
	})
	return z
}
