// Package nsort renders a scene with a Newell–Newell–Sancha depth sort: a
// view-DEPENDENT painter's order recomputed per frame for the current eye.
//
// Position in the layer family:
//   - zsort paints by barycenter depth alone — cheapest, approximate.
//   - bsort builds a view-INDEPENDENT splitting BSP — exact from any eye,
//     but it cuts polygons that straddle partition planes regardless of
//     whether the current view could ever notice, and every cut edge draws
//     as an antialiasing seam. Its build amortises only when world geometry
//     is static and the camera moves.
//   - nsort is exact FOR THE CURRENT EYE and cuts only when an actual
//     occlusion cycle is on screen. On scenes without interpenetration it
//     is a plain depth sort with zero cuts, which makes it the fit for
//     per-frame dynamic geometry (noise fields, mocap) where bsort would
//     rebuild — and re-split — every frame.
package nsort

import (
	"github.com/vibrantgio/seen"
	"github.com/vibrantgio/seen/canvas"
	"github.com/vibrantgio/seen/float"
	"github.com/vibrantgio/seen/layer"
	"github.com/vibrantgio/seen/layer/bsort"
	"github.com/vibrantgio/seen/point"
)

// Layer implements the layer.Layer interface.
type Layer struct {
	scene *seen.Scene
	cache bsort.ShaderCache
	recs  []*record
	stats stats // resolver counters for the last rendered frame
}

var _ layer.Layer = (*Layer)(nil)

func NewLayerForScene(scene *seen.Scene) layer.Layer {
	return &Layer{
		scene: scene,
		cache: make(bsort.ShaderCache),
	}
}

// RenderOn renders all faces of the scene back to front in an order that is
// exact for the current eye (see the package comment).
func (l *Layer) RenderOn(canvas canvas.Canvas) {
	// projection matrix transforms points from world space into camera space
	// and then through viewport prescale and projection matrix into
	// normalized screen space; viewport maps that to real screen space.
	projection := l.scene.Camera.Projection.Mul(l.scene.Viewport.Prescale).Mul(l.scene.Camera.Matrix())
	viewport := l.scene.Viewport.Postscale

	// The shared bsort shader caches world/screen coordinates per face and
	// resolves fill/stroke colors; nsort re-sorts every frame regardless, so
	// only the cache effect matters here, not the change flags.
	shader := bsort.NewShader(l.cache)
	shader.Shade(l.scene, projection, viewport)

	// The eye in world space: the preimage of the eye-space origin under the
	// view transform (see layer/bsort/layer.go for the derivation).
	view := l.scene.Viewport.Prescale.Mul(l.scene.Camera.Matrix())
	var eye point.Point
	if inv, ok := view.Invert(); ok {
		ex, ey, ez := inv.Transform3(0, 0, 0)
		eye = point.Pt(ex, ey, ez)
	} else {
		// Degenerate view (e.g. zero scale). Fall back to the legacy estimate.
		eye = point.Pt(0, 0, -1.0/projection[2][2])
	}

	// Collect the renderable records: world plane from the scene walk,
	// screen projection from the shader cache, culled the same way bsort
	// culls (frustum + backfaces against the world-space eye).
	var planes bsort.Planes
	l.scene.Accept(&planes)
	recs := l.recs[:0]
	for idx := range planes {
		pl := planes[idx]
		frag, present := l.cache[pl.Id]
		if !present || !frag.InViewFrustum {
			continue
		}
		if !l.scene.ShowBackfaces && !frag.Face.ShowBackfaces {
			ed := pl.Normal.Dot(eye)
			pd := pl.Normal.Dot(pl.Barycenter)
			if !float.Equal(ed, pd) && ed < pd {
				continue
			}
		}
		r := &record{
			plane:  pl,
			scr:    frag.ScreenSpace.Points,
			fill:   frag.Fill,
			stroke: frag.Stroke,
			opts:   frag.Face.Options,
		}
		r.computeExtents(eye)
		recs = append(recs, r)
	}

	// project maps a piece's world points through the same clip + viewport
	// path the cached whole-face coordinates went through.
	project := func(pts point.Points) (point.Points, bool) {
		clipped := make(point.Points, len(pts))
		if !pts.Clip(projection, -2, clipped) {
			return nil, false
		}
		screen := make(point.Points, len(clipped))
		clipped.MulB(viewport, screen)
		return screen, true
	}

	l.stats = stats{}
	orderRecords(recs, eye, project, func(r *record) {
		layer.RenderOn(canvas, layer.Fragment{
			Points:  r.scr,
			Fill:    r.fill,
			Stroke:  r.stroke,
			Options: r.opts,
		})
	}, &l.stats)
	l.recs = recs[:0]
}
