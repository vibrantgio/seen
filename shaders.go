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
var FlatShader = Shader(Flat)

func Flat(lights []LightShaderData, surface *SurfaceShaderData, material *Material) color.Color {
	return material.Color
}

// AmbientShader for the `Ambient` shader colors surfaces from ambient light only.
var AmbientShader = Shader(Ambient)

func Ambient(lights []LightShaderData, surface *SurfaceShaderData, material *Material) color.Color {
	c := color.Black
	for _, light := range lights {
		if light.Kind == "ambient" {
			c = applyAmbient(c, light)
		}
	}
	return c.MultiplyChannels(material.Color).Clamp(0.0, 1.0)
}

// DiffusePhong shader implements the Phong shading model with a diffuse
// and ambient term (no specular).
var DiffusePhongShader = Shader(DiffusePhong)

func DiffusePhong(lights []LightShaderData, surface *SurfaceShaderData, material *Material) color.Color {
	c := color.Black
	for _, lsd := range lights {
		switch lsd.Kind {
		case "ambient":
			c = applyAmbient(c, lsd)
		case "directional":
			c = applyDiffuse(c, lsd, lsd.Normal, surface.Normal, material)
		case "point":
			lightNormal := lsd.Point.Subtract(surface.Barycenter).Normalize()
			c = applyDiffuse(c, lsd, lightNormal, surface.Normal, material)
		}
	}
	return c.MultiplyChannels(material.Color).Clamp(0, 1.0)
}

// PhongShader implements the Phong shading model with a diffuse,
// specular, and ambient term.
// See https://en.wikipedia.org/wiki/Phong_reflection_model for more information
var PhongShader = Shader(Phong)

func Phong(lights []LightShaderData, surface *SurfaceShaderData, material *Material) color.Color {
	c := color.Black
	for _, lsd := range lights {
		switch lsd.Kind {
		case "ambient":
			c = applyAmbient(c, lsd)
		case "directional":
			c = applyDiffuseAndSpecular(c, lsd, lsd.Normal, surface.Normal, material)
		case "point":
			lightNormal := lsd.Point.Subtract(surface.Barycenter).Normalize()
			c = applyDiffuseAndSpecular(c, lsd, lightNormal, surface.Normal, material)
		}
	}
	c = c.MultiplyChannels(material.Color)
	if material.Metallic {
		c = c.MinChannels(material.SpecularColor)
	}
	return c.Clamp(0, 1.0)
}

// applyAmbient applies ambient shading
func applyAmbient(c color.Color, lsd LightShaderData) color.Color {
	return c.AddChannels(lsd.Color)
}

// applyDiffuse applies diffuse phong shading
func applyDiffuse(c color.Color, lsd LightShaderData, lightNormal, surfaceNormal Point, material *Material) color.Color {
	dot := lightNormal.Dot(surfaceNormal)
	if dot <= 0.0 {
		return c
	}

	// Apply diffuse phong shading
	return c.AddChannels(lsd.Color.Scale(dot))
}

// applyDiffuseAndSpecular applies diffuse phong shading and specular phong shading.
func applyDiffuseAndSpecular(c color.Color, lsd LightShaderData, lightNormal, surfaceNormal Point, material *Material) color.Color {
	dot := lightNormal.Dot(surfaceNormal)
	if dot <= 0.0 {
		return c
	}

	// Apply diffuse phong shading
	c = c.AddChannels(lsd.Color.Scale(dot))

	// Compute and apply specular phong shading
	reflectionNormal := surfaceNormal.Scale(dot * 2.0).Subtract(lightNormal)
	specularIntensity := math.Pow(0.5+reflectionNormal.Dot(PointZ), material.SpecularExponent) / 255.0
	specularColor := material.SpecularColor.Scale(specularIntensity * lsd.Intensity)
	return c.AddChannels(specularColor)
}
