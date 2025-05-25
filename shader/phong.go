package shader

import (
	"math"

	"github.com/vibrantgio/seen/color"
	"github.com/vibrantgio/seen/intensity"
	"github.com/vibrantgio/seen/light"
	"github.com/vibrantgio/seen/point"
)

// DiffusePhong shader implements the Phong shading model with a diffuse
// and ambient term (no specular).
var DiffusePhong Shader = PhongDiffuseShade

// PhongDiffuseShade applies diffuse phong shading and ignores specular.
func PhongDiffuseShade(lights []light.ShaderData, material *Material, barycenter, normal point.Point) color.Color {
	c := color.Black
	for _, lsd := range lights {
		switch lsd.Kind {
		case light.AmbientKind:
			c = c.AddChannels(lsd.Color)
		case light.DirectionalKind:
			dot := lsd.Normal.Dot(normal)
			if dot > 0.0 {
				c = c.AddChannels(lsd.Color.Scale(dot))
			}
		case light.PointKind:
			dot := lsd.Point.Minus(barycenter).Normalize().Dot(normal)
			if dot > 0.0 {
				c = c.AddChannels(lsd.Color.Scale(dot))
			}
		}
	}
	return c.MultiplyChannels(material.Color).Clamp(0.0, 1.0)
}

// Phong is a short hand for DiffuseAndSpecularPhong.
var Phong Shader = PhongDiffuseAndSpecularShade

// DiffuseAndSpecularPhong implements the Phong shading model with a diffuse,
// specular, and ambient term.
// See https://en.wikipedia.org/wiki/Phong_reflection_model for more information
var DiffuseAndSpecularPhong Shader = PhongDiffuseAndSpecularShade

// PhongDiffuseAndSpecularShade applies diffuse and specular phong shading.
func PhongDiffuseAndSpecularShade(lights []light.ShaderData, material *Material, barycenter, normal point.Point) color.Color {
	apply := func(c color.Color, lsd light.ShaderData, lightNormal, faceNormal point.Point, material *Material) color.Color {
		dot := lightNormal.Dot(faceNormal)
		if dot <= 0.0 {
			return c
		}

		// Apply diffuse phong shading
		c = c.AddChannels(lsd.Color.Scale(dot))

		// Compute and apply specular phong shading
		reflectionNormal := faceNormal.Times(dot * 2.0).Minus(lightNormal)
		specularIntensity := intensity.Intensity(math.Pow(0.5+reflectionNormal.Dot(point.Pt(0, 0, 1)), material.SpecularExponent) / 255.0)
		specularColor := material.SpecularColor.Scale(float64(specularIntensity * lsd.Intensity))
		return c.AddChannels(specularColor)
	}
	c := color.Black
	for _, lsd := range lights {
		switch lsd.Kind {
		case light.AmbientKind:
			c = c.AddChannels(lsd.Color)
		case light.DirectionalKind:
			c = apply(c, lsd, lsd.Normal, normal, material)
		case light.PointKind:
			lightNormal := lsd.Point.Minus(barycenter).Normalize()
			c = apply(c, lsd, lightNormal, normal, material)
		}
	}
	c = c.MultiplyChannels(material.Color)
	if material.Metallic {
		c = c.MinChannels(material.SpecularColor)
	}
	return c.Clamp(0, 1.0)
}
