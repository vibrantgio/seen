package render

import (
	"sort"

	"github.com/reactivego/seen"
)

// SceneLayer extends seen.Scene with a function to paint the
// the scene on a Painter. By implementing this function the
// SceneLayer implements the RenderLayer interface.
type SceneLayer struct {
	*seen.Scene
	renderSurfaces []*RenderSurface
	renderCache    map[string]*RenderSurface
}

func SceneLayerWith(scene *seen.Scene) *SceneLayer {
	return &SceneLayer{
		Scene:          scene,
		renderSurfaces: make([]*RenderSurface, 0, 32),
		renderCache:    make(map[string]*RenderSurface),
	}
}

// Paint creates a RenderSurface for every Surface in the scene's models.
// When encountering a TextShape assign a TextPainter to the RenderSurface.
// When encountering any other shape assign a PathPainter to the RenderSurface.
func (s *SceneLayer) Paint(painter Painter) {
	// projection matrix transforms points from world space into camera space and then
	// through viewport prescale and projection matrix into normalized screen space.
	projection := s.Camera.Projection.Mul(s.Viewport.Prescale).Mul(s.Camera.Matrix())

	// Last transformation from normalized screen space into real screen space.
	viewport := s.Viewport.Postscale

	// Clear out the render models, but reuse the already existing array backing the slice
	s.renderSurfaces = s.renderSurfaces[:0]

	// Process all renderable objects
	s.Model.EachRenderable(func(shape *seen.Shape, lights []seen.LightShaderData, transform seen.Matrix) {
		for _, surface := range shape.Surfaces {
			renderSurface := s.RenderSurfaceWith(&surface, transform, projection, viewport)

			// Assign the correct render function to the render model
			switch shape.Type {
			case "text":
				renderSurface.Render = TextRender
			default:
				renderSurface.Render = PathRender
			}

			// Test projected normal's z-coordinate for culling (if enabled).
			if (s.ShowBackfaces || surface.ShowBackfaces || renderSurface.Normal.Z < 0.0) && renderSurface.InFrustum {
				// Render fill and stroke using material and shader.
				if surface.FillMaterial != nil {
					fill := surface.FillMaterial.Render(lights, s.Shader, renderSurface.ShaderData)
					renderSurface.Fill = &fill
				}
				if surface.StrokeMaterial != nil {
					stroke := surface.StrokeMaterial.Render(lights, s.Shader, renderSurface.ShaderData)
					renderSurface.Stroke = &stroke
				}

				// Round coordinates (if enabled)
				if !s.FractionalPoints {
					pts := renderSurface.ProjectedPoints
					for i, pt := range pts {
						pts[i] = pt.Round()
					}
				}

				// Add the render model to the renderSurfaces slice
				s.renderSurfaces = append(s.renderSurfaces, renderSurface)
			}
		}
	})

	// Sort render models by projected z coordinate. This ensures that the surfaces
	// farthest from the eye are painted first. (Painter's Algorithm)
	sort.Sort(ByZ(s.renderSurfaces))

	// Now for every render model, render it on the given Painter
	for _, rs := range s.renderSurfaces {
		rs.Paint(painter)
	}
}

// RenderSurfaceWith will get or create the renderSurface for
// the given surface. If Regenerate is false, we cache
// these models to reduce object creation and recomputation.
func (s *SceneLayer) RenderSurfaceWith(surface *seen.Surface, transform, projection, viewport seen.Matrix) *RenderSurface {
	if s.Regenerate {
		// No caching
		return RenderSurfaceWith(surface, transform, projection, viewport)
	}
	// Caching enabled, see if its present in the cache
	if rs, ok := s.renderCache[surface.Id]; ok {
		rs.Update(transform, projection, viewport)
		return rs
	}
	// Create new RenderSurface and add to the cache
	rs := RenderSurfaceWith(surface, transform, projection, viewport)
	s.renderCache[surface.Id] = rs
	return rs
}

// FlushCache removes all elements from the cache. This may be necessary
// if you add and remove many shapes from the scene's models since this
// cache has no eviction policy.
func (s *SceneLayer) FlushCache() {
	for k := range s.renderCache {
		delete(s.renderCache, k)
	}
}

// ByZ implements sorting by comparing the Z of the Barycenter of the Projected points
type ByZ []*RenderSurface

func (a ByZ) Len() int {
	return len(a)
}

func (a ByZ) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByZ) Less(i, j int) bool {
	return a[i].Barycenter.Z > a[j].Barycenter.Z
}
