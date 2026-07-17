package bsort

import (
	"slices"

	"github.com/vibrantgio/seen"
	"github.com/vibrantgio/seen/color"
	"github.com/vibrantgio/seen/face"
	"github.com/vibrantgio/seen/light"
	"github.com/vibrantgio/seen/matrix"
	"github.com/vibrantgio/seen/shader"
)

type ShaderFragment struct {
	face.Coordinates
	Fill   *color.Color
	Stroke *color.Color
}

type ShaderCache map[int]ShaderFragment

type Shader struct {
	Cache ShaderCache

	shader     shader.Shader
	projection matrix.Matrix
	viewport   matrix.Matrix
	updated    bool
	world      bool

	lsd struct {
		top   []light.ShaderData
		stack [][]light.ShaderData
	}
}

func NewShader(cache ShaderCache) *Shader {
	return &Shader{Cache: cache}
}

// Shade updates the cache for every face in the scene. updated reports any
// change at all (including projection/viewport-only changes); world reports
// changes to world-space geometry — a face new to the cache, marked Dirty,
// or moved by its model transform. Consumers that derive world-space
// structures (the bsort BSP tree is built from world-space planes) should
// rebuild on world only: a camera or viewport change alone moves no plane.
func (shader *Shader) Shade(scene *seen.Scene, projection, viewport matrix.Matrix) (updated, world bool) {
	shader.shader = scene.Shader
	shader.projection = projection
	shader.viewport = viewport
	shader.updated = false
	shader.world = false
	scene.Accept(shader)
	return shader.updated, shader.world
}

var _ seen.Handler = (*Shader)(nil)

func (v *Shader) EnterGroup() {
	v.lsd.stack = append(v.lsd.stack, slices.Clip(v.lsd.top))
}

func (v *Shader) LeaveGroup() {
	n := len(v.lsd.stack)
	v.lsd.top = v.lsd.stack[n-1]
	v.lsd.stack = v.lsd.stack[:n-1]
}

func (v *Shader) VisitLight(l seen.Light, model matrix.Matrix) {
	if l.IsEnabled() {
		v.lsd.top = append(v.lsd.top, l.ShaderData(model))
	}
}

func (v *Shader) VisitObject(object seen.Object, model matrix.Matrix) {
	for _, f := range object.Faces() {
		frag, present := v.Cache[f.Id]
		if !present {
			frag = ShaderFragment{Coordinates: f.Coordinates(model, v.projection, v.viewport)}
			v.updated = true
			v.world = true
		} else if f.Dirty {
			frag.Coordinates.Update(f.Points, model, v.projection, v.viewport)
			v.updated = true
			v.world = true
		} else {
			// Model must be compared before MaybeUpdate overwrites it: a
			// changed model matrix moves world-space points even when the
			// face's own points are untouched.
			if !model.Equal(frag.Coordinates.Model) {
				v.world = true
			}
			if frag.Coordinates.MaybeUpdate(f.Points, model, v.projection, v.viewport) {
				v.updated = true
			}
		}
		f.Dirty = false

		// Calculate fill and stroke colors using lights, materials and shaders.
		if frag.Coordinates.InViewFrustum {
			barycenter := frag.Coordinates.WorldSpace.Barycenter
			normal := frag.Coordinates.WorldSpace.Normal
			if f.FillMaterial != nil {
				fill := f.FillMaterial.Shade(v.shader, v.lsd.top, barycenter, normal)
				frag.Fill = &fill
			}
			if f.StrokeMaterial != nil {
				stroke := f.StrokeMaterial.Shade(v.shader, v.lsd.top, barycenter, normal)
				frag.Stroke = &stroke
			}
		}
		// Store back unconditionally: Update/MaybeUpdate mutate the local
		// copy's matrices and InViewFrustum flag, and skipping the store for
		// out-of-frustum faces would leave the cached entry permanently
		// stale — re-reporting an update (and re-triggering consumers'
		// rebuilds) every frame thereafter.
		v.Cache[f.Id] = frag
	}
}
