package canvas

import (
	"math"
	"strconv"
	"strings"

	"gioui.org/f32"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"

	"github.com/vibrantgio/seen/canvas"
	"github.com/vibrantgio/seen/color"
	"github.com/vibrantgio/seen/point"
)

// Path
type Path struct {
	*op.Ops
	Points []f32.Point
}

// Set up the path to be painted. Then use Fill and/or Stroke to
// actually paint it.
func (path *Path) Path(points []point.Point) canvas.PathPainter {
	path.Points = nil
	for _, pt := range points {
		path.Points = append(path.Points, f32.Pt(float32(pt.X), float32(pt.Y)))
	}
	return path
}

// Fill the path
func (p *Path) Fill(style canvas.Style) {
	if len(p.Points) == 0 {
		return
	}
	fillOpacity := 1.0
	if o, present := style["fill-opacity"]; present {
		if o, err := strconv.ParseFloat(o, 64); err == nil {
			fillOpacity = math.Min(math.Max(0.0, o), 1.0)
		}
	}
	if c, present := style["fill"]; present {
		if fill, err := color.ColorWithString(c); err == nil {
			fill.A = fillOpacity
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
	state := clip.Outline{Path: path.End()}.Op().Push(p.Ops)
	paint.PaintOp{}.Add(p.Ops)
	state.Pop()
}

// Stroke the outline of the path.
// Key "stroke-width" is supported in style.
func (p *Path) Stroke(style canvas.Style) {
	if len(p.Points) == 0 {
		return
	}
	strokeWidth := 1
	if sw, present := style["stroke-width"]; present {
		sw = strings.TrimSuffix(sw, "px")
		if sw, err := strconv.Atoi(sw); err == nil {
			strokeWidth = sw
		}
	}
	if c, present := style["stroke"]; present {
		if stroke, err := color.ColorWithString(c); err == nil {
			paint.ColorOp{Color: stroke.NRGBA()}.Add(p.Ops)
		}
	}
	var path clip.Path
	path.Begin(p.Ops)
	path.MoveTo(p.Points[0])
	for _, p := range p.Points[1:] {
		path.LineTo(p)
	}
	path.Close()
	state := clip.Stroke{Path: path.End(), Width: float32(strokeWidth)}.Op().Push(p.Ops)
	paint.PaintOp{}.Add(p.Ops)
	state.Pop()
}
