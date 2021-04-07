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
	renderModels     []*RenderModel
	renderModelCache map[string]*RenderModel
}

func MakeSceneLayer(scene *seen.Scene) *SceneLayer {
	s := &SceneLayer{Scene: scene}
	s.Init()
	return s
}

func (s *SceneLayer) Init() {
	s.renderModels = make([]*RenderModel, 0, 32)
	s.renderModelCache = make(map[string]*RenderModel)
}

// Paint creates a RenderModel for every Surface in the scene's models.
// When encountering a TextShape assign a TextPainter to the RenderModel.
// When encountering any other shape assign a PathPainter to the RenderModel.
func (s *SceneLayer) Paint(painter Painter) {
	// projection matrix transforms points from world space into camera space and then
	// trhough viewport prescale and projection matrix into normalized screen space.
	projection := s.Camera.Projection.Mul(s.Viewport.Prescale).Mul(s.Camera.Matrix())

	// Last transformation from normalized screen space into real screen space.
	viewport := s.Viewport.Postscale

	// Clear out the render models, but reuse the already existing array backing the slice
	s.renderModels = s.renderModels[:0]

	// Process all renderable objects
	s.Model.EachRenderable(
		// ShapeFunc
		func(shape *seen.Shape, lights []seen.LightShaderData, transform seen.Matrix) {
			for _, surface := range shape.Surfaces {
				renderModel := s.makeRenderModel(&surface, transform, projection, viewport)

				// Assign the correct render function to the render model
				switch shape.Kind {
				case "text":
					renderModel.Render = TextRender
				default:
					renderModel.Render = PathRender
				}

				// Test projected normal's z-coordinate for culling (if enabled).
				if (s.ShowBackfaces || surface.ShowBackfaces || renderModel.Normal.Z < 0.0) && renderModel.InFrustrum {
					// Render fill and stroke using material and shader.
					if surface.FillMaterial != nil {
						fill := surface.FillMaterial.Render(lights, s.Shader, renderModel.ShaderData)
						renderModel.Fill = &fill
					}
					if surface.StrokeMaterial != nil {
						stroke := surface.StrokeMaterial.Render(lights, s.Shader, renderModel.ShaderData)
						renderModel.Stroke = &stroke
					}

					// Round coordinates (if enabled)
					if !s.FractionalPoints {
						pts := renderModel.ProjectedPoints
						for i, pt := range pts {
							pts[i] = pt.Round()
						}
					}

					// Add the render model to the renderModels slice
					s.renderModels = append(s.renderModels, renderModel)
				}
			}
		})

	// Sort render models by projected z coordinate. This ensures that the surfaces
	// farthest from the eye are painted first. (Painter's Algorithm)
	sort.Sort(ByZ(s.renderModels))

	// Now for every render model, render it on the given Painter
	for _, m := range s.renderModels {
		m.Paint(painter)
	}
}

// makeRenderModel will get or create the rendermodel for
// the given surface. If Regenerate is false, we cache
// these models to reduce object creation and recomputation.
func (s *SceneLayer) makeRenderModel(surface *seen.Surface, transform, projection, viewport seen.Matrix) *RenderModel {
	if s.Regenerate {
		// No caching
		return MakeRenderModel(surface, transform, projection, viewport)
	}
	// Caching enabled, see if its present in the cache
	m, ok := s.renderModelCache[surface.Id]
	if ok {
		m.Update(transform, projection, viewport)
		return m
	}
	// Create new RenderModel and add to the cache
	m = MakeRenderModel(surface, transform, projection, viewport)
	s.renderModelCache[surface.Id] = m
	return m
}

// FlushCache removes all elements from the cache. This may be necessary
// if you add and remove many shapes from the scene's models since this
// cache has no eviction policy.
func (s *SceneLayer) FlushCache() {
	s.renderModelCache = make(map[string]*RenderModel)
}

// ByZ implements sorting by comparing the Z of the Barycenter of the Projected points
type ByZ []*RenderModel

func (a ByZ) Len() int {
	return len(a)
}

func (a ByZ) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByZ) Less(i, j int) bool {
	return a[i].Barycenter.Z > a[j].Barycenter.Z
}
