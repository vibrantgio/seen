package layer

import (
	"strconv"

	"github.com/vibrantgio/seen/affine"
	"github.com/vibrantgio/seen/canvas"
)

func RenderOn(canvas canvas.Canvas, fragment Fragment) {
	options := fragment.Options
	points := fragment.Points
	fill := fragment.Fill
	stroke := fragment.Stroke

	// Special case is rendering text...
	if text, present := options["text"]; present {
		style := map[string]string{
			"fill":        "none",
			"text-anchor": "middle",
		}
		if fill != nil {
			style["fill"] = fill.Hex()
		}
		if font, present := options["font"]; present {
			style["font"] = font
		}
		if family, present := options["font-family"]; present {
			style["font-family"] = family
		}
		if size, present := options["font-size"]; present {
			style["font-size"] = size
		}
		if weight, present := options["font-weight"]; present {
			style["font-weight"] = weight
		}
		if anchor, present := options["anchor"]; present {
			style["text-anchor"] = anchor
		}
		if length, present := options["inline-size"]; present {
			style["inline-size"] = length
		}
		xform := affine.SolveForAffineTransform(affine.Basis(points))
		canvas.Text().FillText(xform, text, style)
		return
	}

	// Default case is to render a path.
	path := canvas.Path()
	path.Path(points)

	if fill != nil {
		style := map[string]string{
			"fill":         fill.Hex(),
			"fill-opacity": strconv.FormatFloat(fill.A, 'f', -1, 64),
		}
		path.Fill(style)
	}

	if stroke != nil {
		style := map[string]string{
			"fill":         "none",
			"stroke":       stroke.Hex(),
			"text-anchor":  "middle",
			"stroke-width": "1",
		}
		if v, present := options["stroke-width"]; present {
			style["stroke-width"] = v
		}
		path.Stroke(style)
	}
}
