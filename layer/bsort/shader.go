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

	lsd struct {
		top   []light.ShaderData
		stack [][]light.ShaderData
	}
}

func NewShader(cache ShaderCache) *Shader {
	return &Shader{Cache: cache}
}

func (shader *Shader) Shade(scene *seen.Scene, projection, viewport matrix.Matrix) bool {
	shader.shader = scene.Shader
	shader.projection = projection
	shader.viewport = viewport
	shader.updated = false
	scene.Accept(shader)
	return shader.updated
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
			v.Cache[f.Id] = frag
			v.updated = true
		} else {
			if f.Dirty {
				frag.Coordinates.Update(f.Points, model, v.projection, v.viewport)
				v.updated = true
			} else if frag.Coordinates.MaybeUpdate(f.Points, model, v.projection, v.viewport) {
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
				v.Cache[f.Id] = frag
			}
			if f.StrokeMaterial != nil {
				stroke := f.StrokeMaterial.Shade(v.shader, v.lsd.top, barycenter, normal)
				frag.Stroke = &stroke
				v.Cache[f.Id] = frag
			}
		}
	}
}
