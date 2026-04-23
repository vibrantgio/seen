package canvas

import (
	"image"

	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"

	"github.com/vibrantgio/seen/canvas"
	"github.com/vibrantgio/seen/color"
)

// Rect
type Rect struct {
	*op.Ops
	Width  float32
	Height float32
	Rx     float32
	Ry     float32
}

func (rect *Rect) Rect(width, height float64) canvas.RectPainter {
	rect.Width, rect.Height = float32(width), float32(height)
	return rect
}

func (rect *Rect) CornerRadius(rx, ry float64) canvas.RectPainter {
	rect.Rx, rect.Ry = float32(rx), float32(ry)
	return rect
}

// Fill the rect
func (rect *Rect) Fill(style canvas.Style) {
	if c, present := style["fill"]; present {
		if fill, err := color.ColorWithString(c); err == nil {
			paint.ColorOp{Color: fill.NRGBA()}.Add(rect.Ops)
		}
	}
	if rect.Rx == 0.0 && rect.Ry == 0.0 {
		state := clip.Rect(image.Rect(0, 0, int(rect.Width), int(rect.Height))).Push(rect.Ops)
		paint.PaintOp{}.Add(rect.Ops)
		state.Pop()
	} else if rect.Rx == rect.Ry {
		state := clip.UniformRRect(image.Rect(0, 0, int(rect.Width), int(rect.Height)), int(rect.Rx)).Push(rect.Ops)
		paint.PaintOp{}.Add(rect.Ops)
		state.Pop()
	} else {

		// TBD
	}

}
