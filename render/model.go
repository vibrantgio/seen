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

	Transform *seen.Matrix

	Projection *seen.Matrix

	Viewport *seen.Matrix

	Points []seen.Point

	ShaderData *seen.SurfaceShaderData

	ProjectedPoints []seen.Point
	Barycenter      *seen.Point
	Normal          *seen.Point

	InFrustrum bool

	Fill *colors.Color

	Stroke *colors.Color
}

func MakeRenderModel(surface *seen.Surface, transform, projection, viewport *seen.Matrix) *RenderModel {
	m := &RenderModel{}
	m.Init(surface, transform, projection, viewport)
	return m
}

func (m *RenderModel) Init(surface *seen.Surface, transform, projection, viewport *seen.Matrix) {
	m.Surface = surface
	m.Transform = transform
	m.Projection = projection
	m.Viewport = viewport
	m.Points = surface.Points
	m.update()
}

func (m *RenderModel) Paint(painter Painter) {
	m.Render(m, painter)
}

func (m *RenderModel) Update(transform, projection, viewport *seen.Matrix) {
	if m.Surface.Dirty || !transform.Equal(m.Transform) || !projection.Equal(m.Projection) || !viewport.Equal(m.Viewport) {
		m.Transform = transform
		m.Projection = projection
		m.Viewport = viewport
		m.update()
	}
}

func (m *RenderModel) update() {

	// Apply model transform to surface points. Calculates transformed points and barycenter
	worldSpacePoints, baryCenter := m.Transform.TransformPoints(m.Points)
	// Initialize the shader data with the baryCenter and the normal of the transformed points.
	m.ShaderData = &seen.SurfaceShaderData{baryCenter, seen.MakePointNormal(worldSpacePoints)}

	// Transform into camera space and check whether points are inside the frustrum along the way.
	var cameraSpaceCoords = make([]seen.Coordinate, len(worldSpacePoints))
	m.InFrustrum = true
	for i := range cameraSpaceCoords {
		c := m.Projection.TransformCoordinate(worldSpacePoints[i].ToCoordinate())
		if c.Z <= -2 {
			m.InFrustrum = false
		}
		cameraSpaceCoords[i] = *c
	}

	// Project camera space points into screen space
	m.ProjectedPoints, m.Barycenter = m.Viewport.ProjectCoordinatesToPoints(cameraSpaceCoords)
	// Compute the surface normal in screen space.
	m.Normal = seen.MakePointNormal(m.ProjectedPoints)

	// Surface has been updated, we can clear the Dirty flag
	m.Surface.Dirty = false
}
