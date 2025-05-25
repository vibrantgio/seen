package zsort

import (
	"cmp"
	"slices"

	"github.com/vibrantgio/seen"
	"github.com/vibrantgio/seen/canvas"
	"github.com/vibrantgio/seen/face"
	"github.com/vibrantgio/seen/layer"
	"github.com/vibrantgio/seen/light"
	"github.com/vibrantgio/seen/matrix"
)

// Layer implements the layer.Layer interface.
// Layer extends seen.Scene with a function to paint the
// the scene on a Painter.
type Layer struct {
	scene      *seen.Scene
	cache      map[int]face.Coordinates
	zfragments []zfragment
}

type zfragment struct {
	z        float64
	fragment layer.Fragment
}

var _ layer.Layer = (*Layer)(nil)

func NewLayerForScene(scene *seen.Scene) layer.Layer {
	return &Layer{
		scene:      scene,
		cache:      make(map[int]face.Coordinates),
		zfragments: make([]zfragment, 0, 32),
	}
}

// RenderOn renders every Face in the scene's groups on the Canvas passed in as an argument.
// When encountering a TextShape use a TextPainter to render.
// When encountering any other shape assign a PathPainter to render.
func (l *Layer) RenderOn(canvas canvas.Canvas) {
	// projection matrix transforms points from world space into camera space and then
	// through viewport prescale and projection matrix into normalized screen space.
	projection := l.scene.Camera.Projection.Mul(l.scene.Viewport.Prescale).Mul(l.scene.Camera.Matrix())

	// Last transformation from normalized screen space into real screen space.
	viewport := l.scene.Viewport.Postscale

	// Clear out the fragments, but reuse the already existing array backing the slice
	l.zfragments = l.zfragments[:0]

	// Process all renderable objects
	l.scene.Group.EachRenderable(func(object seen.Object, lights []light.ShaderData, model matrix.Matrix) {
		faces := object.Faces()
		for i := range faces {
			f := &faces[i]

			// Get or create the coordinates for the given face.
			var coordinates face.Coordinates
			if l.scene.Regenerate {
				// No caching
				coordinates = f.Coordinates(model, projection, viewport)
				if !l.scene.FractionalPoints {
					coordinates.ScreenSpace.Points.Round()
				}
			} else {
				var updated, present bool
				if coordinates, present = l.cache[f.Id]; !present {
					coordinates = f.Coordinates(model, projection, viewport)
					updated = true
				} else {
					if f.Dirty {
						coordinates.Update(f.Points, model, projection, viewport)
						updated = true
					} else {
						updated = coordinates.MaybeUpdate(f.Points, model, projection, viewport)
					}
				}
				if updated {
					if !l.scene.FractionalPoints {
						coordinates.ScreenSpace.Points.Round()
					}
					l.cache[f.Id] = coordinates
				}
			}
			f.Dirty = false

			// Test projected normal's z-coordinate for culling (if enabled).
			if (l.scene.ShowBackfaces || f.ShowBackfaces || coordinates.ScreenSpace.Normal.Z < 0.0) && coordinates.InViewFrustum {
				var zfrag zfragment
				zfrag.z = coordinates.ScreenSpace.Barycenter.Z
				zfrag.fragment.Points = coordinates.ScreenSpace.Points

				// Render fill and stroke using material and shader.
				barycenter := coordinates.WorldSpace.Barycenter
				normal := coordinates.WorldSpace.Normal
				if f.FillMaterial != nil {
					fill := f.FillMaterial.Shade(l.scene.Shader, lights, barycenter, normal)
					zfrag.fragment.Fill = &fill
				}
				if f.StrokeMaterial != nil {
					stroke := f.StrokeMaterial.Shade(l.scene.Shader, lights, barycenter, normal)
					zfrag.fragment.Stroke = &stroke
				}
				zfrag.fragment.Options = f.Options

				// Insert fragments by projected z coordinate. This ensures that
				// the fragments for faces farthest from the eye are painted
				// first. (Painter's Algorithm)
				index, _ := slices.BinarySearchFunc(l.zfragments, zfrag, func(e, t zfragment) int {
					return -cmp.Compare(e.z, t.z)
				})
				l.zfragments = slices.Insert(l.zfragments, index, zfrag)
			}
		}
	})

	// Now for every fragment, render its layer.Fragment on the given canvas.
	// The faces farthest from the eye are painted first. (Painter's Algorithm)
	for _, zf := range l.zfragments {
		layer.RenderOn(canvas, zf.fragment)
	}
}

// FlushCache removes all elements from the cache. This may be necessary
// if you add and remove many shapes from the scene's groups since this
// cache has no eviction policy.
func (s *Layer) FlushCache() {
	for k := range s.cache {
		delete(s.cache, k)
	}
}

// ByZ implements sorting by comparing the Z of the Barycenter of the Projected points
type ByZ []zfragment

func (a ByZ) Len() int {
	return len(a)
}

func (a ByZ) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByZ) Less(i, j int) bool {
	return a[i].z < a[j].z
}
