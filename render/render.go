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
	fill := "none"
	if model.Fill != nil {
		fill = model.Fill.Hex()
	}
	font, fontPresent := model.Surface.Options["font"]
	if !fontPresent {
		font = ""
	}
	anchor, anchorPresent := model.Surface.Options["anchor"]
	if !anchorPresent {
		anchor = "middle"
	}
	style := map[string]string{
		"fill":        fill,
		"font":        font,
		"text-anchor": anchor,
	}
	xform := affine.SolveForAffineTransform(model.ProjectedPoints)
	text := model.Surface.Options["text"]
	painter.Text().FillText(xform, text, style)
}
