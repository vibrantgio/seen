package bsp

import (
	"github.com/reactivego/seen"
	"github.com/reactivego/seen/render"
)

type SurfaceCache map[int]*render.Surface

func (cache SurfaceCache) Update(scene *seen.Scene, projection, viewport seen.Matrix) (updated bool) {
	updatecache := &UpdateCache{Scene: scene, Projection: projection, Viewport: viewport, Cache: cache}
	scene.Group.Accept(updatecache)
	return updatecache.Updated
}

type UpdateCache struct {
	TransformStack
	Stack [][]seen.LightShaderData
	LSD   []seen.LightShaderData

	Scene      *seen.Scene
	Projection seen.Matrix
	Viewport   seen.Matrix
	Cache      SurfaceCache
	Updated    bool
}

func (v *UpdateCache) Push() {
	v.TransformStack.Push()
	v.Stack = append(v.Stack, v.LSD)
}

func (v *UpdateCache) Pop() {
	v.LSD = v.Stack[len(v.Stack)-1]
	v.Stack = v.Stack[:len(v.Stack)-1]
	v.TransformStack.Pop()
}

func (v *UpdateCache) VisitLight(l *seen.Light) {
	if l.Enabled {
		v.LSD = append(v.LSD, l.ShaderData(v.Transform.Mul(l.Matrix())))
	}
}

func (v *UpdateCache) VisitSurface(surface *seen.Surface) {
	var rs *render.Surface
	if cs, ok := v.Cache[surface.Id]; ok {
		if cs.Update(v.Transform, v.Projection, v.Viewport) {
			v.Updated = true
		}
		rs = cs
	} else {
		rs = render.SurfaceWith(surface, v.Transform, v.Projection, v.Viewport)
		v.Cache[surface.Id] = rs
		v.Updated = true
	}

	// Test projected normal's z-coordinate for culling (if enabled).
	// ShowBackfaces := (v.Scene.ShowBackfaces || surface.ShowBackfaces || rs.Normal.Z < 0.0)
	ShowBackfaces := true
	if ShowBackfaces && rs.InFrustum {
		// Render fill and stroke using material and shader.
		if surface.FillMaterial != nil {
			fill := surface.FillMaterial.Render(v.LSD, v.Scene.Shader, rs.ShaderData)
			rs.Fill = &fill
		}
		if surface.StrokeMaterial != nil {
			stroke := surface.StrokeMaterial.Render(v.LSD, v.Scene.Shader, rs.ShaderData)
			rs.Stroke = &stroke
		}
	}
}
