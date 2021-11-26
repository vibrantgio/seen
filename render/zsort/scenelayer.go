package zsort

import (
	"sort"

	"github.com/reactivego/seen"
	"github.com/reactivego/seen/render"
)

// SceneLayer extends seen.Scene with a function to paint the
// the scene on a Painter. By implementing this function the
// SceneLayer implements the render.Layer interface.
type SceneLayer struct {
	*seen.Scene
	surfaces []*render.RenderSurface
	cache    map[int]*render.RenderSurface
}

func LayerWith(scene *seen.Scene) *SceneLayer {
	return &SceneLayer{
		Scene:    scene,
		surfaces: make([]*render.RenderSurface, 0, 32),
		cache:    make(map[int]*render.RenderSurface),
	}
}

// Paint creates a RenderSurface for every Surface in the scene's groups.
// When encountering a TextShape assign a TextPainter to the RenderSurface.
// When encountering any other shape assign a PathPainter to the RenderSurface.
func (s *SceneLayer) Paint(painter render.Painter) {
	// projection matrix transforms points from world space into camera space and then
	// through viewport prescale and projection matrix into normalized screen space.
	projection := s.Camera.Projection.Mul(s.Viewport.Prescale).Mul(s.Camera.Matrix())

	// Last transformation from normalized screen space into real screen space.
	viewport := s.Viewport.Postscale

	// Clear out the render surfaces, but reuse the already existing array backing the slice
	s.surfaces = s.surfaces[:0]

	// Process all renderable objects
	s.Group.EachRenderable(func(shape *seen.Shape, lights []seen.LightShaderData, transform seen.Matrix) {
		for _, surface := range shape.Surfaces {
			// Get or create the renderSurface for the given surface.
			var rs *render.RenderSurface
			// If Regenerate is false, we cache the render surfaces to reduce object
			// creation and recomputation.
			if s.Regenerate {
				// No caching
				surface.Shape = shape
				rs = render.RenderSurfaceWith(&surface, transform, projection, viewport)
			} else {
				// Caching enabled, see if its present in the cache
				if cs, ok := s.cache[surface.Id]; ok {
					cs.Update(transform, projection, viewport)
					rs = cs
				}
				// Create new RenderSurface and add to the cache
				surface.Shape = shape
				rs = render.RenderSurfaceWith(&surface, transform, projection, viewport)
				s.cache[surface.Id] = rs
			}
			// Test projected normal's z-coordinate for culling (if enabled).
			if (s.ShowBackfaces || surface.ShowBackfaces || rs.Normal.Z < 0.0) && rs.InFrustum {
				// Render fill and stroke using material and shader.
				if surface.FillMaterial != nil {
					fill := surface.FillMaterial.Render(lights, s.Shader, rs.ShaderData)
					rs.Fill = &fill
				}
				if surface.StrokeMaterial != nil {
					stroke := surface.StrokeMaterial.Render(lights, s.Shader, rs.ShaderData)
					rs.Stroke = &stroke
				}

				// Round coordinates (if enabled)
				if !s.FractionalPoints {
					pts := rs.ProjectedPoints
					for i, pt := range pts {
						pts[i] = pt.Round()
					}
				}

				// Add the render surface to the surfaces slice
				s.surfaces = append(s.surfaces, rs)
			}
		}
	})

	// Sort render surfaces by projected z coordinate. This ensures that the surfaces
	// farthest from the eye are painted first. (Painter's Algorithm)
	sort.Sort(sort.Reverse(ByZ(s.surfaces)))

	// Now for every render surface, render it on the given Painter
	for _, rs := range s.surfaces {
		rs.Paint(painter)
	}
}

// FlushCache removes all elements from the cache. This may be necessary
// if you add and remove many shapes from the scene's groups since this
// cache has no eviction policy.
func (s *SceneLayer) FlushCache() {
	for k := range s.cache {
		delete(s.cache, k)
	}
}

// ByZ implements sorting by comparing the Z of the Barycenter of the Projected points
type ByZ []*render.RenderSurface

func (a ByZ) Len() int {
	return len(a)
}

func (a ByZ) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByZ) Less(i, j int) bool {
	return a[i].Barycenter.Z < a[j].Barycenter.Z
}
