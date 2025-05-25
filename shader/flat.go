package shader

import (
	"github.com/vibrantgio/seen/color"
	"github.com/vibrantgio/seen/light"
	"github.com/vibrantgio/seen/point"
)

// Flat for the `Flat` shader colors faces with the material color, disregarding all
// light sources.
var Flat Shader = FlatShade

func FlatShade(lights []light.ShaderData, material *Material, barycenter, normal point.Point) color.Color {
	return material.Color
}
