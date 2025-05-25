package shader

import (
	"github.com/vibrantgio/seen/color"
	"github.com/vibrantgio/seen/light"
	"github.com/vibrantgio/seen/point"
)

// Shader implements the Shade method
type Shader func(lights []light.ShaderData, material *Material, barycenter, normal point.Point) color.Color

// Shade
// `lights` is an object containing the ambient, point, and directional light sources.
// `face` is an instance of `FaceData` and contains the transformed and projected face data.
// `material` is an instance of `Material` and contains the color and other attributes for determining how light reflects off the face.
func (shade Shader) Shade(lights []light.ShaderData, material *Material, barycenter, normal point.Point) color.Color {
	return shade(lights, material, barycenter, normal)
}

var Default Shader = Phong
