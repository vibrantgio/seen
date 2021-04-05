package render

import "github.com/reactivego/seen/affine"

// PathRender renders a RenderModel for a Path surface onto a Painter.
func PathRender(model *RenderModel, painter Painter) {
	path := painter.Path()
	path.Path(model.ProjectedPoints)

	if model.Fill != nil {
		path.Fill(map[string]string{
			"fill":         model.Fill.Hex(),
			"fill-opacity": Ftoa(model.Fill.A),
		})
	}

	if model.Stroke != nil {
		strokeWidth := "1"
		if v, ok := model.Surface.Options["stroke-width"]; ok {
			strokeWidth = v
		}
		path.Stroke(map[string]string{
			"fill":         "none",
			"stroke":       model.Stroke.Hex(),
			"stroke-width": strokeWidth,
		})
	}
}

// TextRender renders a RenderModel for a Text Surface onto a Painter.
func TextRender(model *RenderModel, painter Painter) {
	style := map[string]string{
		"fill":        "none",
		"text-anchor": "middle",
	}
	if model.Fill != nil {
		style["fill"] = model.Fill.Hex()
	}
	if font, present := model.Surface.Options["font"]; present {
		style["font"] = font
	}
	if family, present := model.Surface.Options["font-family"]; present {
		style["font-family"] = family
	}
	if size, present := model.Surface.Options["font-size"]; present {
		style["font-size"] = size
	}
	if weight, present := model.Surface.Options["font-weight"]; present {
		style["font-weight"] = weight
	}
	if anchor, present := model.Surface.Options["anchor"]; present {
		style["text-anchor"] = anchor
	}
	xform := affine.SolveForAffineTransform(model.ProjectedPoints)
	text := model.Surface.Options["text"]
	painter.Text().FillText(xform, text, style)
}
