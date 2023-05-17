package seen

import (
	"math"

	"github.com/reactivego/seen/color"
)

// Shader implements the Shade method
type Shader func(lights []LightShaderData, surface *SurfaceShaderData, material *Material) color.Color

// Shade
// `lights` is an object containing the ambient, point, and directional light sources.
// `surface` is an instance of `SurfaceShaderData` and contains the transformed and projected surface data.
// `material` is an instance of `Material` and contains the color and other attributes for determining how light reflects off the surface.
func (shade Shader) Shade(lights []LightShaderData, surface *SurfaceShaderData, material *Material) color.Color {
	return shade(lights, surface, material)
}

// SurfaceShaderData
type SurfaceShaderData struct {
	Barycenter Point
	Normal     Point
}

// Shade for the `Flat` shader colors surfaces with the material color, disregarding all
// light sources.
var FlatShader Shader = FlatShade

func FlatShade(lights []LightShaderData, surface *SurfaceShaderData, material *Material) color.Color {
	return material.Color
}

// AmbientShader for the `Ambient` shader colors surfaces from ambient light only.
var AmbientShader Shader = AmbientShade

func AmbientShade(lights []LightShaderData, surface *SurfaceShaderData, material *Material) color.Color {
	c := color.Black
	for _, lsd := range lights {
		if lsd.Kind == AmbientKind {
			c = c.AddChannels(lsd.Color)
		}
	}
	return c.MultiplyChannels(material.Color).Clamp(0.0, 1.0)
}

// DiffusePhong shader implements the Phong shading model with a diffuse
// and ambient term (no specular).
var DiffusePhongShader Shader = PhongDiffuseShade

// PhongDiffuseShade applies diffuse phong shading and ignores specular.
func PhongDiffuseShade(lights []LightShaderData, surface *SurfaceShaderData, material *Material) color.Color {
	c := color.Black
	for _, lsd := range lights {
		switch lsd.Kind {
		case AmbientKind:
			c = c.AddChannels(lsd.Color)
		case DirectionalKind:
			dot := lsd.Normal.Dot(surface.Normal)
			if dot > 0.0 {
				c = c.AddChannels(lsd.Color.Scale(dot))
			}
		case PointKind:
			dot := lsd.Point.Subtract(surface.Barycenter).Normalize().Dot(surface.Normal)
			if dot > 0.0 {
				c = c.AddChannels(lsd.Color.Scale(dot))
			}
		}
	}
	return c.MultiplyChannels(material.Color).Clamp(0.0, 1.0)
}

// PhongShader implements the Phong shading model with a diffuse,
// specular, and ambient term.
// See https://en.wikipedia.org/wiki/Phong_reflection_model for more information
var PhongShader Shader = PhongDiffuseAndSpecularShade

// PhongDiffuseAndSpecularShade applies diffuse and specular phong shading.
func PhongDiffuseAndSpecularShade(lights []LightShaderData, surface *SurfaceShaderData, material *Material) color.Color {
	apply := func(c color.Color, lsd LightShaderData, lightNormal, surfaceNormal Point, material *Material) color.Color {
		dot := lightNormal.Dot(surfaceNormal)
		if dot <= 0.0 {
			return c
		}

		// Apply diffuse phong shading
		c = c.AddChannels(lsd.Color.Scale(dot))

		// Compute and apply specular phong shading
		reflectionNormal := surfaceNormal.Scale(dot * 2.0).Subtract(lightNormal)
		specularIntensity := math.Pow(0.5+reflectionNormal.Dot(Point{0, 0, 1}), material.SpecularExponent) / 255.0
		specularColor := material.SpecularColor.Scale(specularIntensity * lsd.Intensity)
		return c.AddChannels(specularColor)
	}
	c := color.Black
	for _, lsd := range lights {
		switch lsd.Kind {
		case AmbientKind:
			c = c.AddChannels(lsd.Color)
		case DirectionalKind:
			c = apply(c, lsd, lsd.Normal, surface.Normal, material)
		case PointKind:
			lightNormal := lsd.Point.Subtract(surface.Barycenter).Normalize()
			c = apply(c, lsd, lightNormal, surface.Normal, material)
		}
	}
	c = c.MultiplyChannels(material.Color)
	if material.Metallic {
		c = c.MinChannels(material.SpecularColor)
	}
	return c.Clamp(0, 1.0)
}

var DefaultShader Shader = PhongDiffuseAndSpecularShade
