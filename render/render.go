package render

import "github.com/reactivego/seen/affine"

// PathRender renders a RenderSurface for a Path surface onto a Painter.
func PathRender(surface *RenderSurface, painter Painter) {
	path := painter.Path()
	path.Path(surface.ProjectedPoints)

	if surface.Fill != nil {
		path.Fill(map[string]string{
			"fill":         surface.Fill.Hex(),
			"fill-opacity": Ftoa(surface.Fill.A),
		})
	}

	if surface.Stroke != nil {
		strokeWidth := "1"
		if v, ok := surface.Surface.Options["stroke-width"]; ok {
			strokeWidth = v
		}
		path.Stroke(map[string]string{
			"fill":         "none",
			"stroke":       surface.Stroke.Hex(),
			"stroke-width": strokeWidth,
		})
	}
}

// TextRender renders a RenderSurface for a Text Surface onto a Painter.
func TextRender(surface *RenderSurface, painter Painter) {
	style := map[string]string{
		"fill":        "none",
		"text-anchor": "middle",
	}
	if surface.Fill != nil {
		style["fill"] = surface.Fill.Hex()
	}
	if font, present := surface.Surface.Options["font"]; present {
		style["font"] = font
	}
	if family, present := surface.Surface.Options["font-family"]; present {
		style["font-family"] = family
	}
	if size, present := surface.Surface.Options["font-size"]; present {
		style["font-size"] = size
	}
	if weight, present := surface.Surface.Options["font-weight"]; present {
		style["font-weight"] = weight
	}
	if anchor, present := surface.Surface.Options["anchor"]; present {
		style["text-anchor"] = anchor
	}
	if length, present := surface.Surface.Options["textLength"]; present {
		style["textLength"] = length
	}
	xform := affine.SolveForAffineTransform(surface.ProjectedPoints)
	text := surface.Surface.Options["text"]
	painter.Text().FillText(xform, text, style)
}
