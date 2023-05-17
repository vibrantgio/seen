package bsp

import (
	"github.com/reactivego/seen"
	"github.com/reactivego/seen/float"
	"github.com/reactivego/seen/render"
)

// Layer implements the render.Layer interface.
// Layer extends seen.Scene with a function to paint the
// the scene on a Painter.
type Layer struct {
	*seen.Scene

	// BSP is a Binary Space Partitioning generated for the scene.
	// The BSP must be re-generated whenever the scene graph geometry is modified.
	bsp *BSP

	surfaces []*render.Surface
	cache    SurfaceCache
}

var _ render.Layer = (*Layer)(nil)

func NewLayerForScene(scene *seen.Scene) render.Layer {
	return &Layer{
		Scene:    scene,
		surfaces: make([]*render.Surface, 0, 32),
		cache:    make(map[int]*render.Surface),
	}
}

// Paint creates a render.Surface for every seen.Surface in the scene's objects.
// When encountering a TextShape assign a TextPainter to the render.Surface.
// When encountering any other shape assign a PathPainter to the render.Surface.
func (s *Layer) Paint(painter render.Painter) {
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
	s.bsp.Display(eye, func(plane []Plane) {
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
