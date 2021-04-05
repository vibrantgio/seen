package gio

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	"eliasnaur.com/font/roboto/robotobold"
	"eliasnaur.com/font/roboto/robotoregular"

	"golang.org/x/image/math/fixed"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/font/opentype"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"

	"github.com/reactivego/seen"
	"github.com/reactivego/seen/affine"
	"github.com/reactivego/seen/colors"
	"github.com/reactivego/seen/render"
)

var roboto struct {
	sync.Once
	faces []text.FontFace
}

func RobotoFontFaces() []text.FontFace {
	register := func(fnt text.Font, ttf []byte) {
		face, err := opentype.Parse(ttf)
		if err != nil {
			panic(fmt.Sprintf("failed to parse font: %v", err))
		}
		fnt.Typeface = "Roboto"
		roboto.faces = append(roboto.faces, text.FontFace{Font: fnt, Face: face})
	}
	roboto.Do(func() {
		// Weight: Normal (400)
		register(text.Font{Style: text.Regular, Weight: text.Normal}, robotoregular.TTF)
		// Weight: Bold (600)
		register(text.Font{Style: text.Regular, Weight: text.Bold}, robotobold.TTF)
	})
	return roboto.faces
}

var (
	RobotoNormal = text.Font{
		Typeface: "Roboto",
		Variant:  "",
		Style:    text.Regular,
		Weight:   text.Normal /*400*/}
	RobotoBold = text.Font{
		Typeface: "Roboto",
		Variant:  "",
		Style:    text.Regular,
		Weight:   text.Bold /*600*/}

	shaper = text.NewCache(RobotoFontFaces())
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

func (c *Context) Draw(ops *op.Ops, queue event.Queue) {
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
	if seen.Scheduler.Run() {
		op.InvalidateOp{}.Add(ops)
	}
}

func (c *Context) Animate() seen.Animator {
	animator := seen.MakeAnimator()
	animator.OnFrame(func(d, dt time.Duration) {
		c.Render()
	})
	return animator
}

func (c *Context) Drag(options ...seen.DragOption) *seen.Drag {
	drag := seen.MakeDrag(options...)
	c.inputs = append(c.inputs, func(ops *op.Ops) {
		defer op.Save(ops).Load()
		pointer.PassOp{Pass: true}.Add(ops)
		const types = pointer.Press | pointer.Drag | pointer.Release
		pointer.InputOp{Tag: drag, Types: types}.Add(ops)
	})
	previous := struct {
		Position f32.Point
		Time     time.Duration
	}{}
	c.handlers = append(c.handlers, func(q event.Queue) {
		for _, event := range q.Events(drag) {
			if p, ok := event.(pointer.Event); ok {
				switch p.Type {
				case pointer.Press:
					drag.Handle(seen.DragEvent{
						Type: seen.DragStart,
						X:    float64(p.Position.X),
						Y:    float64(p.Position.Y),
						T:    p.Time,
					})
				case pointer.Drag:
					if previous.Time != 0 {
						dP := p.Position.Sub(previous.Position)
						dT := p.Time - previous.Time
						drag.Handle(seen.DragEvent{
							Type: seen.DragMove,
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
						drag.Handle(seen.DragEvent{
							Type: seen.DragEnd,
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
	return drag
}

func (c *Context) Zoom() *seen.Zoom {
	zoom := seen.MakeZoom()
	c.inputs = append(c.inputs, func(ops *op.Ops) {
		defer op.Save(ops).Load()
		pointer.PassOp{Pass: true}.Add(ops)
		pointer.InputOp{
			Tag:          zoom,
			Types:        pointer.Scroll,
			ScrollBounds: image.Rect(-120, -120, 120, 120),
		}.Add(ops)
	})
	c.handlers = append(c.handlers, func(q event.Queue) {
		for _, event := range q.Events(zoom) {
			if p, ok := event.(pointer.Event); ok {
				dx, dy := -float64(p.Scroll.X), -float64(p.Scroll.Y)
				dxy := math.Copysign(math.Hypot(dx, dy), dy)
				zoom.Handle(seen.ZoomEvent{
					Type:       seen.ZoomMove,
					Dx:         dx,
					Dy:         dy,
					WheelDelta: dxy,
				})
			}
		}
	})
	return zoom
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
func (p *TextPainter) FillText(t affine.Matrix, txt string, style render.Style) {
	defer op.Save(p.Ops).Load()
	aff := f32.NewAffine2D(float32(t.A), float32(t.C), float32(t.E), float32(t.B), float32(t.D), float32(t.F))
	op.Affine(aff).Add(p.Ops)

	font := RobotoNormal
	size := 10
	fill := color.NRGBA{0, 0, 0, 255}

	if family, present := style["font-family"]; present {
		font.Typeface = text.Typeface(family)
	}
	if weight, present := style["font-weight"]; present {
		switch weight {
		case "normal":
			font.Weight = text.Normal
		case "bold":
			font.Weight = text.Bold
		}
	}
	if sz, present := style["font-size"]; present {
		if strings.HasSuffix(sz, "px") {
			sz = sz[:len(sz)-2]
		}
		if sz, err := strconv.Atoi(sz); err == nil {
			size = sz
		}
	}
	if c, present := style["fill"]; present {
		if f, err := colors.ColorWithString(c); err == nil {
			fill = f.NRGBA()
		}
	}
	ax, ay := float32(0.5), float32(1.0)
	if a, present := style["text-anchor"]; present {
		if a == "start" {
			ax = 0.0
		} else if a == "end" {
			ax = 1.0
		}
	}
	textLength := 2000
	if tl, present := style["textLength"]; present {
		if strings.HasSuffix(tl, "px") {
			tl = tl[:len(tl)-2]
		}
		if tl, err := strconv.Atoi(tl); err == nil {
			textLength = tl
		}
	}

	// Layout the txt string given font and size and a MaxLineWidth
	lines := shaper.LayoutString(font, fixed.I(size), textLength, txt)

	// Determine the size of the layout rectangle dx,dy
	dx, dy := float32(0), float32(0)
	for _, line := range lines {
		dy += float32(line.Ascent.Ceil() + line.Descent.Ceil())
		lineWidth := float32(line.Width.Ceil())
		if dx < lineWidth {
			dx = lineWidth
		}
	}

	// Actually paint the txt
	offset := f32.Pt(-ax*dx, -ay*dy)
	for _, line := range lines {
		state := op.Save(p.Ops)
		offset.Y += float32(line.Ascent.Ceil())
		op.Offset(offset).Add(p.Ops)
		offset.Y += float32(line.Descent.Ceil())
		shaper.Shape(font, fixed.I(size), line.Layout).Add(p.Ops)
		paint.ColorOp{Color: fill}.Add(p.Ops)
		paint.PaintOp{}.Add(p.Ops)
		state.Load()
	}
}
