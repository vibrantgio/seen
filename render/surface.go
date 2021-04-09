package render

import (
	"github.com/reactivego/seen"
	"github.com/reactivego/seen/colors"
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
	// Render is a reference to a specific render function to be used to render
	// the RenderSurface on a Painter.
	Render func(*RenderSurface, Painter)

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
	m := &RenderSurface{}
	m.Surface = surface
	m.Points = surface.Points

	m.Transform = transform
	m.Projection = projection
	m.Viewport = viewport
	m.update()
	return m
}

func (m *RenderSurface) Paint(painter Painter) {
	m.Render(m, painter)
}

func (m *RenderSurface) Update(transform, projection, viewport seen.Matrix) {
	if m.Surface.Dirty || !transform.Equal(m.Transform) || !projection.Equal(m.Projection) || !viewport.Equal(m.Viewport) {
		m.Transform = transform
		m.Projection = projection
		m.Viewport = viewport
		m.update()
	}
}

func (m *RenderSurface) update() {
	if len(m.WorldSpacePoints) != len(m.Points) {
		m.WorldSpacePoints = make([]seen.Point, len(m.Points))
	}
	if len(m.ProjectedPoints) != len(m.Points) {
		m.ProjectedPoints = make([]seen.Point, len(m.Points))
	}

	// Apply model transform to surface points. Calculates transformed points and barycenter
	wsBaryCenter := m.Points.Mul(m.Transform, m.WorldSpacePoints)
	wsNormal := m.WorldSpacePoints.Normal().Normalize()

	// Initialize the shader data with the baryCenter and the normal of the transformed points.
	m.ShaderData = &seen.SurfaceShaderData{Barycenter: wsBaryCenter, Normal: wsNormal}

	var clippedPoints = make(seen.Points, len(m.WorldSpacePoints))
	if m.InFrustum = m.WorldSpacePoints.Clip(m.Projection, -2, clippedPoints); m.InFrustum {
		// Project camera space points into screen space
		m.Barycenter = clippedPoints.Mul(m.Viewport, m.ProjectedPoints)
		m.Normal = m.ProjectedPoints.Normal().Normalize()

		// Surface has been updated, we can clear the Dirty flag
		m.Surface.Dirty = false
	}
}
