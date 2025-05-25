package shader

import (
	"github.com/vibrantgio/seen/color"
	"github.com/vibrantgio/seen/light"
	"github.com/vibrantgio/seen/point"
)

// Ambient for the `Ambient` shader colors faces from ambient light only.
var Ambient Shader = AmbientShade

func AmbientShade(lights []light.ShaderData, material *Material, barycenter, normal point.Point) color.Color {
	c := color.Black
	for _, lsd := range lights {
		if lsd.Kind == light.AmbientKind {
			c = c.AddChannels(lsd.Color)
		}
	}
	return c.MultiplyChannels(material.Color).Clamp(0.0, 1.0)
}
