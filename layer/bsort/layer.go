package bsort

import (
	"github.com/vibrantgio/seen"
	"github.com/vibrantgio/seen/canvas"
	"github.com/vibrantgio/seen/float"
	"github.com/vibrantgio/seen/layer"
	"github.com/vibrantgio/seen/layer/bsort/bsp"
	"github.com/vibrantgio/seen/point"
)

// Layer implements the layer.Layer interface.
type Layer struct {
	scene *seen.Scene
	cache ShaderCache
	tree  *bsp.Tree
}

var _ layer.Layer = (*Layer)(nil)

func NewLayerForScene(scene *seen.Scene) layer.Layer {
	return &Layer{
		scene: scene,
		cache: make(ShaderCache),
	}
}

// RenderOn renders all faces that are part of the shapes and objects
// in the scene on to the given canvas
func (s *Layer) RenderOn(canvas canvas.Canvas) {
	// projection matrix transforms points from world space into camera space and then
	// through viewport prescale and projection matrix into normalized screen space.
	projection := s.scene.Camera.Projection.Mul(s.scene.Viewport.Prescale).Mul(s.scene.Camera.Matrix())

	// Last transformation from normalized screen space into real screen space.
	viewport := s.scene.Viewport.Postscale

	shader := NewShader(s.cache)
	if shader.Shade(s.scene, projection, viewport) || s.tree == nil {
		var planes Planes
		s.scene.Accept(&planes)
		s.tree = bsp.NewTree(planes)
		// fmt.Printf("#planes %d\n", len(collector.Planes))
	}

	// Find out where the eye is located.
	eye := point.Pt(0, 0, -1.0/projection[2][2])
	// fmt.Printf("eye: %v\n", eye)

	// Walk the bsp tree and render the render face back to front
	s.tree.Display(eye, func(plane []bsp.Plane) {
		for i := range plane {
			frag := s.cache[plane[i].Id]
			if frag.InViewFrustum {
				if !s.scene.ShowBackfaces && !frag.Face.ShowBackfaces {
					ed := plane[i].Normal.Dot(eye)
					pd := plane[i].Normal.Dot(plane[i].Barycenter)
					if !float.Equal(ed, pd) && ed < pd {
						continue
					}
				}
				fragment := layer.Fragment{
					Points:  frag.ScreenSpace.Points,
					Fill:    frag.Fill,
					Stroke:  frag.Stroke,
					Options: frag.Face.Options,
				}
				layer.RenderOn(canvas, fragment)
			}
		}
	})
}
