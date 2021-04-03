package gio

import (
	"image"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"

	"github.com/reactivego/seen"
	"github.com/reactivego/seen/affine"
	"github.com/reactivego/seen/colors"
	"github.com/reactivego/seen/render"
)

// Context
type Context struct {
	window   *app.Window
	render   []func(render.Painter)
	inputs   []func(*op.Ops)
	handlers []func(event.Queue)
}

// MakeContext creates a render context for the given op.Ops and layer.
func MakeContext(window *app.Window, layer render.RenderLayer) *Context {
	context := &Context{window: window}
	if layer != nil {
		context.Layer(layer)
	}
	return context
}

func (c *Context) Render() {
	c.window.Invalidate()
}

func (c *Context) HandleEvents(queue event.Queue, ops *op.Ops) {
	p := Painter{Ops: ops}
	for _, render := range c.render {
		render(&p)
	}
	for _, input := range c.inputs {
		input(ops)
	}
	for _, handler := range c.handlers {
		handler(queue)
	}
}

func (c *Context) Animate() seen.Animator {
	animator := seen.MakeAnimator()
	animator.OnFrame(func(d, dt float64) {
		c.Render()
	})
	return animator
}

func (c *Context) Drag(options ...seen.DragOption) *seen.Drag {
	drag := seen.MakeDrag(options...)
	c.inputs = append(c.inputs, func(ops *op.Ops) {
		defer op.Save(ops).Load()
		pointer.InputOp{Tag: drag, Types: pointer.Drag}.Add(ops)
	})
	c.handlers = append(c.handlers, func(q event.Queue) {
		pos := f32.Pt(0, 0)
		for _, event := range q.Events(drag) {
			if p, ok := event.(pointer.Event); ok && p.Type == pointer.Drag {
				if pos != f32.Pt(0, 0) {
					d := p.Position.Sub(pos)
					drag.Handle(seen.DragEvent{
						OffsetRelativeX: float64(d.X),
						OffsetRelativeY: float64(d.Y),
					})
				}
				pos = p.Position
			}
		}
	})
	return drag
}

func (c *Context) Layer(layer render.RenderLayer) {
	c.render = append(c.render, layer.Paint)
}

type Painter struct {
	*op.Ops
	PathPainter
	RectPainter
	CirclePainter
	TextPainter
}

func (c *Painter) Path() render.PathPainter {
	c.PathPainter.Ops = c.Ops
	return &c.PathPainter
}

func (c *Painter) Rect() render.RectPainter {
	c.RectPainter.Ops = c.Ops
	return &c.RectPainter
}

func (c *Painter) Circle() render.CirclePainter {
	c.CirclePainter.Ops = c.Ops
	return &c.CirclePainter
}

func (c *Painter) Text() render.TextPainter {
	c.TextPainter.Ops = c.Ops
	return &c.TextPainter
}

func (c *Painter) Reset() {

}

func (c *Painter) Cleanup() {

}

// PathPainter
type PathPainter struct {
	*op.Ops
	Points []f32.Point
}

// Set up the path to be painted. Then use Fill and/or Stroke to
// actually paint it.
func (p *PathPainter) Path(points []seen.Point) {
	p.Points = nil
	for _, pt := range points {
		p.Points = append(p.Points, f32.Pt(float32(pt.X), float32(pt.Y)))
	}
}

// Fill the path
func (p *PathPainter) Fill(style render.Style) {
	if len(p.Points) == 0 {
		return
	}
	defer op.Save(p.Ops).Load()
	if c, present := style["fill"]; present {
		if fill, err := colors.ColorWithString(c); err == nil {
			paint.ColorOp{Color: fill.NRGBA()}.Add(p.Ops)
		}
	}
	var path clip.Path
	path.Begin(p.Ops)
	path.MoveTo(p.Points[0])
	for _, p := range p.Points[1:] {
		path.LineTo(p)
	}
	path.Close()
	clip.Outline{Path: path.End()}.Op().Add(p.Ops)
	paint.PaintOp{}.Add(p.Ops)
}

// Stroke the outline of the path.
// Key "stroke-width" is supported in style.
func (p *PathPainter) Stroke(render.Style) {

}

// RectPainter
type RectPainter struct {
	*op.Ops
	Width  float32
	Height float32
	Rx     float32
	Ry     float32
}

func (p *RectPainter) Size(width, height float64) {
	p.Width, p.Height = float32(width), float32(height)
}

func (p *RectPainter) CornerRadius(rx, ry float64) {
	p.Rx, p.Ry = float32(rx), float32(ry)
}

// Fill the rect
func (p *RectPainter) Fill(style render.Style) {
	defer op.Save(p.Ops).Load()
	if c, present := style["fill"]; present {
		if fill, err := colors.ColorWithString(c); err == nil {
			paint.ColorOp{Color: fill.NRGBA()}.Add(p.Ops)
		}
	}
	if p.Rx == 0.0 && p.Ry == 0.0 {
		clip.Rect(image.Rect(0, 0, int(p.Width), int(p.Height))).Add(p.Ops)
	} else if p.Rx == p.Ry {
		clip.UniformRRect(f32.Rect(0, 0, p.Width, p.Height), p.Rx).Add(p.Ops)
	} else {
		// TBD
	}
	paint.PaintOp{}.Add(p.Ops)
}

// CirclePainter
type CirclePainter struct{ *op.Ops }

func (p *CirclePainter) Fill(render.Style) {
}

// TextPainter
type TextPainter struct{ *op.Ops }

// FillText
// transform is an affine matrix approximating a 3D transform of the plane on which the text is to be painted.
// text is the text to be painted.
// Style supports the following keys: fill, font, text-anchor
func (p *TextPainter) FillText(transform affine.Matrix, text string, style render.Style) {
}
