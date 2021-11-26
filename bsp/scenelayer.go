package bsp

import (
	"github.com/reactivego/seen"
	"github.com/reactivego/seen/float"
	"github.com/reactivego/seen/render"
)

// SceneLayer extends seen.Scene with a function to paint the
// the scene on a Painter. By implementing this function the
// SceneLayer implements the RenderLayer interface.
type SceneLayer struct {
	*seen.Scene

	// BSP is a Binary Space Partitioning generated for the scene.
	// The BSP must be re-generated whenever the scene graph geometry is modified.
	bsp *BSP

	surfaces []*render.RenderSurface
	cache    SurfaceCache
}

func SceneLayerWith(scene *seen.Scene) *SceneLayer {
	return &SceneLayer{
		Scene:    scene,
		surfaces: make([]*render.RenderSurface, 0, 32),
		cache:    make(map[int]*render.RenderSurface),
	}
}

// Paint creates a RenderSurface for every Surface in the scene's objects.
// When encountering a TextShape assign a TextPainter to the RenderSurface.
// When encountering any other shape assign a PathPainter to the RenderSurface.
func (s *SceneLayer) Paint(painter render.Painter) {
	// projection matrix transforms points from world space into camera space and then
	// through viewport prescale and projection matrix into normalized screen space.
	projection := s.Scene.Camera.Projection.Mul(s.Scene.Viewport.Prescale).Mul(s.Scene.Camera.Matrix())

	// Last transformation from normalized screen space into real screen space.
	viewport := s.Scene.Viewport.Postscale

	// Update render surfaces in the cache
	if s.cache.Update(s.Scene, projection, viewport) || s.bsp == nil {
		buildbsp := &Builder{}
		s.Scene.Group.Accept(buildbsp)
		s.bsp = buildbsp.Build()
		// fmt.Printf("#planes %d\n", len(buildbsp.Planes))
	}

	// Find out where the eye is located.
	eye := seen.Pt(0, 0, -1.0/projection[2][2])
	// fmt.Printf("eye: %v\n", eye)

	// Walk the bsp tree and render the render surface back to front
	s.bsp.Display(eye, func(plane []seen.Plane) {
		for i := range plane {
			rs := s.cache[plane[i].Id]
			if rs.InFrustum {
				if !s.Scene.ShowBackfaces && !rs.Surface.ShowBackfaces {
					ed := plane[i].Normal.Dot(eye)
					pd := plane[i].Normal.Dot(plane[i].Barycenter)
					if !float.Equal(ed, pd) && ed < pd {
						continue
					}
				}
				rs.Paint(painter)
			}
		}
	})
}
