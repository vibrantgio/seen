package render

import (
	"github.com/reactivego/seen"
	"github.com/reactivego/seen/colors"
)

// RenderModel contains the transformed and projected points as
// well as various data needed to shade and paint a `Surface`.
//
// Once initialized, the object will have a constant memory footprint down to
// `Number` primitives. Also, we compare each transform and projection to
// prevent unnecessary re-computation.
//
// If you need to force a re-computation, mark the surface as 'dirty'.
//
// RenderModel manages the painting of a single Surface.
type RenderModel struct {
	// Render is a reference to a specific render function to be used to render
	// the RenderModel on a Painter.
	Render func(*RenderModel, Painter)

	// Surface is a reference to the Surface that is being painted.
	// The reference is retained so it can be checked for the Dirty flag.
	// When the Dirty flag is set, the RenderModel needs to be regenerated.
	Surface *seen.Surface
	Points  []seen.Point

	Transform seen.Matrix

	Projection seen.Matrix

	Viewport seen.Matrix

	ShaderData       *seen.SurfaceShaderData
	WorldSpacePoints []seen.Point
	ProjectedPoints  []seen.Point
	Barycenter       seen.Point
	Normal           seen.Point

	InFrustrum bool

	Fill *colors.Color

	Stroke *colors.Color
}

func MakeRenderModel(surface *seen.Surface, transform, projection, viewport seen.Matrix) *RenderModel {
	m := &RenderModel{}
	m.Init(surface, transform, projection, viewport)
	return m
}

func (m *RenderModel) Init(surface *seen.Surface, transform, projection, viewport seen.Matrix) {
	m.Surface = surface
	m.Points = surface.Points

	m.Transform = transform
	m.Projection = projection
	m.Viewport = viewport
	m.update()
}

func (m *RenderModel) Paint(painter Painter) {
	m.Render(m, painter)
}

func (m *RenderModel) Update(transform, projection, viewport seen.Matrix) {
	if m.Surface.Dirty || !transform.Equal(m.Transform) || !projection.Equal(m.Projection) || !viewport.Equal(m.Viewport) {
		m.Transform = transform
		m.Projection = projection
		m.Viewport = viewport
		m.update()
	}
}

func (m *RenderModel) update() {
	if len(m.WorldSpacePoints) != len(m.Points) {
		m.WorldSpacePoints = make([]seen.Point, len(m.Points))
	}
	if len(m.ProjectedPoints) != len(m.Points) {
		m.ProjectedPoints = make([]seen.Point, len(m.Points))
	}

	// Apply model transform to surface points. Calculates transformed points and barycenter
	wsBaryCenter := m.Transform.TransformPoints(m.Points, m.WorldSpacePoints)
	wsNormal := seen.PointNormal(m.WorldSpacePoints).Normalize()
	// Initialize the shader data with the baryCenter and the normal of the transformed points.
	m.ShaderData = &seen.SurfaceShaderData{Barycenter: wsBaryCenter, Normal: wsNormal}

	// Transform into camera space and check whether points are inside the frustrum along the way.
	var cameraSpaceCoords = make([]seen.Coordinate, len(m.WorldSpacePoints))
	m.InFrustrum = true
	for i := range cameraSpaceCoords {
		c := m.Projection.TransformCoordinate(m.WorldSpacePoints[i].ToCoordinate())
		if c.Z <= -2 {
			m.InFrustrum = false
		}
		cameraSpaceCoords[i] = c
	}

	// Project camera space points into screen space
	m.Barycenter = m.Viewport.ProjectCoordinatesToPoints(cameraSpaceCoords, m.ProjectedPoints)
	// Compute the surface normal in screen space.
	m.Normal = seen.PointNormal(m.ProjectedPoints).Normalize()

	// Surface has been updated, we can clear the Dirty flag
	m.Surface.Dirty = false
}
