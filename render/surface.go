package render

import (
	"github.com/reactivego/seen"
	"github.com/reactivego/seen/affine"
	colors "github.com/reactivego/seen/color"
)

// RenderSurface contains the transformed and projected points as
// well as various data needed to shade and paint a `Surface`.
//
// Once initialized, the object will have a constant memory footprint down to
// `Number` primitives. Also, we compare each transform and projection to
// prevent unnecessary re-computation.
//
// If you need to force a re-computation, mark the surface as 'dirty'.
//
// RenderSurface manages the painting of a single Surface.
type RenderSurface struct {
	// Paint is the paint function to be used to paint this surface.
	Paint func(Painter)

	// Surface is a reference to the Surface that is being painted.
	// The reference is retained so it can be checked for the Dirty flag.
	// When the Dirty flag is set, the RenderSurface needs to be regenerated.
	Surface *seen.Surface
	Points  seen.Points

	Transform seen.Matrix

	Projection seen.Matrix

	Viewport seen.Matrix

	ShaderData       *seen.SurfaceShaderData
	WorldSpacePoints seen.Points
	ProjectedPoints  seen.Points
	Barycenter       seen.Point
	Normal           seen.Point

	InFrustum bool

	Fill *colors.Color

	Stroke *colors.Color
}

func RenderSurfaceWith(surface *seen.Surface, transform, projection, viewport seen.Matrix) *RenderSurface {
	rs := &RenderSurface{}
	// Assign the correct render function to the render model
	if surface.Shape.Type == "text" {
		rs.Paint = rs.PaintText
	} else {
		rs.Paint = rs.PaintPath
	}
	rs.Surface = surface
	rs.Points = surface.Points

	rs.Transform = transform
	rs.Projection = projection
	rs.Viewport = viewport
	rs.update()
	return rs
}

func (rs *RenderSurface) Update(transform, projection, viewport seen.Matrix) (updated bool) {
	if rs.Surface.Dirty || !transform.Equal(rs.Transform) || !projection.Equal(rs.Projection) || !viewport.Equal(rs.Viewport) {
		rs.Transform = transform
		rs.Projection = projection
		rs.Viewport = viewport
		rs.update()
		updated = true
	}
	return
}

func (rs *RenderSurface) update() {
	if len(rs.WorldSpacePoints) != len(rs.Points) {
		rs.WorldSpacePoints = make([]seen.Point, len(rs.Points))
	}
	if len(rs.ProjectedPoints) != len(rs.Points) {
		rs.ProjectedPoints = make([]seen.Point, len(rs.Points))
	}

	// Apply model transform to surface points. Calculates transformed points and barycenter
	wsBaryCenter := rs.Points.Mul(rs.Transform, rs.WorldSpacePoints)
	wsNormal := rs.WorldSpacePoints.Normal().Normalize()

	// Initialize the shader data with the baryCenter and the normal of the transformed points.
	rs.ShaderData = &seen.SurfaceShaderData{Barycenter: wsBaryCenter, Normal: wsNormal}

	var clippedPoints = make(seen.Points, len(rs.WorldSpacePoints))
	if rs.InFrustum = rs.WorldSpacePoints.Clip(rs.Projection, -2, clippedPoints); rs.InFrustum {
		// Project camera space points into screen space
		rs.Barycenter = clippedPoints.Mul(rs.Viewport, rs.ProjectedPoints)
		rs.Normal = rs.ProjectedPoints.Normal().Normalize()

		// Surface has been updated, we can clear the Dirty flag
		rs.Surface.Dirty = false
	}
}

// PaintPath paints a path render surface onto a Painter.
func (rs *RenderSurface) PaintPath(painter Painter) {
	path := painter.Path()
	path.Path(rs.ProjectedPoints)

	if rs.Fill != nil {
		path.Fill(map[string]string{
			"fill":         rs.Fill.Hex(),
			"fill-opacity": Ftoa(rs.Fill.A),
		})
	}

	if rs.Stroke != nil {
		strokeWidth := "1"
		if v, ok := rs.Surface.Options["stroke-width"]; ok {
			strokeWidth = v
		}
		path.Stroke(map[string]string{
			"fill":         "none",
			"stroke":       rs.Stroke.Hex(),
			"stroke-width": strokeWidth,
		})
	}
}

// PaintText paints a text render surface onto a Painter.
func (rs *RenderSurface) PaintText(painter Painter) {
	style := map[string]string{
		"fill":        "none",
		"text-anchor": "middle",
	}
	if rs.Fill != nil {
		style["fill"] = rs.Fill.Hex()
	}
	if font, present := rs.Surface.Options["font"]; present {
		style["font"] = font
	}
	if family, present := rs.Surface.Options["font-family"]; present {
		style["font-family"] = family
	}
	if size, present := rs.Surface.Options["font-size"]; present {
		style["font-size"] = size
	}
	if weight, present := rs.Surface.Options["font-weight"]; present {
		style["font-weight"] = weight
	}
	if anchor, present := rs.Surface.Options["anchor"]; present {
		style["text-anchor"] = anchor
	}
	if length, present := rs.Surface.Options["textLength"]; present {
		style["textLength"] = length
	}
	xform := affine.SolveForAffineTransform(rs.ProjectedPoints)
	text := rs.Surface.Options["text"]
	painter.Text().FillText(xform, text, style)
}
