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
	scene   *seen.Scene
	cache   ShaderCache
	tree    *bsp.Tree
	noSplit bool
}

var _ layer.Layer = (*Layer)(nil)

func NewLayerForScene(scene *seen.Scene) layer.Layer {
	return &Layer{
		scene: scene,
		cache: make(ShaderCache),
	}
}

// NewNoSplitLayerForScene returns the bsort layer in whole-polygon mode: the
// BSP keeps every straddling face whole on its barycenter's side instead of
// cutting it (bsp.NewTreeNoSplit). Ordering is then approximate for
// interpenetrating or cyclically occluding faces, but no cut edges are ever
// introduced — each cut edge in the splitting mode renders as a visible
// antialiasing seam across the face's fill, and on animated geometry the
// tree (and so the seam pattern) is rebuilt every frame, making the seams
// crawl. Prefer this mode for scenes whose depth order cannot cycle from any
// eye position (height fields, non-intersecting meshes); use
// NewLayerForScene when faces genuinely interpenetrate and exact order is
// worth the seams.
func NewNoSplitLayerForScene(scene *seen.Scene) layer.Layer {
	return &Layer{
		scene:   scene,
		cache:   make(ShaderCache),
		noSplit: true,
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
		if s.noSplit {
			s.tree = bsp.NewTreeNoSplit(planes)
		} else {
			s.tree = bsp.NewTree(planes)
		}
		// fmt.Printf("#planes %d\n", len(collector.Planes))
	}

	// Find out where the eye (center of projection) is located in world space.
	//
	// The BSP planes, their barycenters and normals are all in world space, so
	// Display needs the eye in world space too. The center of projection is the
	// world point that maps to the eye-space origin under the view transform
	// Prescale * Camera.Matrix() (the projection matrix's row 3 of [0,0,-1,0]
	// makes w_clip vanish exactly there). It is therefore the preimage of the
	// origin under that affine view transform, i.e. (Prescale * Camera)^-1 * 0.
	//
	// This is independent of the frustum's near/far and correctly accounts for
	// camera dolly and viewport offset. The previous formula
	//   point.Pt(0, 0, -1.0/projection[2][2])
	// ignored all translations and was additionally off by a (f-n)/(f+n) factor
	// even for an identity camera with a symmetric viewport.
	view := s.scene.Viewport.Prescale.Mul(s.scene.Camera.Matrix())
	var eye point.Point
	if inv, ok := view.Invert(); ok {
		ex, ey, ez := inv.Transform3(0, 0, 0)
		eye = point.Pt(ex, ey, ez)
	} else {
		// Degenerate view (e.g. zero scale). Fall back to the legacy estimate.
		eye = point.Pt(0, 0, -1.0/projection[2][2])
	}
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
