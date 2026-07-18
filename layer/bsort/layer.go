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

// NewLayerForScene returns the splitting-BSP layer: exact painter's order
// from ANY eye, with polygons that straddle a partition plane cut into
// pieces at build time. The tree is world-space and rebuilt only when world
// geometry changes, so its build cost amortises across camera motion —
// which makes this the layer for STATIC geometry under a moving eye. For
// per-frame dynamic geometry (noise fields, mocap) use layer/nsort, which
// orders for the current eye and cuts only on genuine occlusion cycles.
func NewLayerForScene(scene *seen.Scene) layer.Layer {
	return &Layer{
		scene: scene,
		cache: make(ShaderCache),
	}
}

// RenderOn renders all faces that are part of the shapes and objects
// in the scene on to the given canvas
func (s *Layer) RenderOn(canvas canvas.Canvas) {
	// projection matrix transforms points from world space through the
	// camera's view (world transform, eye, normalization) and projection
	// into normalized screen space.
	projection := s.scene.Camera.Projection.Mul(s.scene.Camera.View())

	// Last transformation from normalized screen space into real screen space.
	viewport := s.scene.Viewport.Screen

	shader := NewShader(s.cache)
	// The tree is built from world-space planes, so it only goes stale when
	// world geometry changes; camera/viewport-only changes reuse it — the
	// eye passed to Display below adapts the traversal to the new view.
	if _, world := shader.Shade(s.scene, projection, viewport); world || s.tree == nil {
		var planes Planes
		s.scene.Accept(&planes)
		s.tree = bsp.NewTree(planes)
		// fmt.Printf("#planes %d\n", len(collector.Planes))
	}

	// The eye (center of projection) in world space. The BSP planes, their
	// barycenters and normals are all in world space, so Display needs the
	// eye in world space too (see camera.EyeInWorld for the derivation).
	eye := s.scene.Camera.EyeInWorld()

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
				points := frag.ScreenSpace.Points
				if plane[i].Piece {
					// A piece split off a straddling face renders its own
					// polygon — the cached coordinates hold the whole face.
					// Project the piece's world-space points through the
					// same clip + viewport path the cache went through.
					clipped := make(point.Points, len(plane[i].Points))
					if !plane[i].Points.Clip(projection, -2, clipped) {
						continue
					}
					screen := make(point.Points, len(clipped))
					clipped.MulB(viewport, screen)
					points = screen
				}
				fragment := layer.Fragment{
					Points:  points,
					Fill:    frag.Fill,
					Stroke:  frag.Stroke,
					Options: frag.Face.Options,
				}
				layer.RenderOn(canvas, fragment)
			}
		}
	})
}
