package seen

import (
	"math"
	"github.com/reactivego/seen/colors"
)

// Shader implements the Shade method
type Shader interface {
	// Shade
	// `lights` is an object containing the ambient, point, and directional light sources.
	// `surface` is an instance of `SurfaceShaderData` and contains the transformed and projected surface data.
	// `material` is an instance of `Material` and contains the color and other attributes for determining how light reflects off the surface.
	Shade(lights []*LightRenderData, surface *SurfaceShaderData, material *Material) *colors.Color
}

// SurfaceShaderData
type SurfaceShaderData struct {
	Barycenter *Point
	Normal     *Point
}

// FlatShader
type FlatShader struct {
}

func MakeFlatShader() Shader {
	return &FlatShader{}
}

// Shade for the `Flat` shader colors surfaces with the material color, disregarding all
// light sources.
func (s *FlatShader) Shade(lights []*LightRenderData, surface *SurfaceShaderData, material *Material) *colors.Color {
	return material.Color
}

// AmbientShader
type AmbientShader struct {
}

func MakeAmbientShader() Shader {
	return &AmbientShader{}
}

// Shade for the `Ambient` shader colors surfaces from ambient light only.
func (s *AmbientShader) Shade(lights []*LightRenderData, surface *SurfaceShaderData, material *Material) *colors.Color {
	c := colors.Black
	for _, light := range lights {
		if light.Kind == "ambient" {
			c = applyAmbient(c, light)
		}
	}
	return c.MultiplyChannels(material.Color).Clamp(0.0, 1.0)
}

// DiffusePhong shader implements the Phong shading model with a diffuse
// and ambient term (no specular).
type DiffusePhongShader struct {
}

func MakeDiffusePhongShader() Shader {
	return &DiffusePhongShader{}
}

func (s *DiffusePhongShader) Shade(lights []*LightRenderData, surface *SurfaceShaderData, material *Material) *colors.Color {
	c := colors.Black
	for _, light := range lights {
		switch light.Kind {
		case "ambient":
			c = applyAmbient(c, light)
		case "directional":
			c = applyDiffuse(c, light, light.Normal, surface.Normal, material)
		case "point":
			lightNormal := light.Point.Subtract(surface.Barycenter).Normalize()
			c  = applyDiffuse(c, light, lightNormal, surface.Normal, material)
		}
	}
	return c.MultiplyChannels(material.Color).Clamp(0, 1.0)
}

// PhongShader implements the Phong shading model with a diffuse,
// specular, and ambient term.
// See https://en.wikipedia.org/wiki/Phong_reflection_model for more information
type PhongShader struct {
}

// MakePhongShader
func MakePhongShader() Shader {
	return &PhongShader{}
}

// Shade
func (s *PhongShader) Shade(lights []*LightRenderData, surface *SurfaceShaderData, material *Material) *colors.Color {
	c := colors.Black
	for _,light := range lights {
		switch light.Kind {
		case "ambient":
			c = applyAmbient(c, light)
		case "directional":
			c = applyDiffuseAndSpecular(c, light, light.Normal, surface.Normal, material)
		case "point":
			lightNormal := light.Point.Subtract(surface.Barycenter).Normalize()
			c = applyDiffuseAndSpecular(c, light, lightNormal, surface.Normal, material)
		}
	}
	c = c.MultiplyChannels(material.Color)
	if material.Metallic {
		c = c.MinChannels(material.SpecularColor)
	}
	return c.Clamp(0, 1.0)
}

// applyAmbient applies ambient shading
func applyAmbient(c *colors.Color, light *LightRenderData) *colors.Color {
	return c.AddChannels(light.ColorIntensity)
}

// applyDiffuse applies diffuse phong shading
func applyDiffuse(c *colors.Color, light *LightRenderData, lightNormal, surfaceNormal *Point, material *Material) *colors.Color {
	dot := lightNormal.Dot(surfaceNormal)
	if dot <= 0.0 {
		return c
	}

	// Apply diffuse phong shading
	return c.AddChannels(light.ColorIntensity.Scale(dot))
}

// applyDiffuseAndSpecular applies diffuse phong shading and specular phong shading.
func applyDiffuseAndSpecular(c *colors.Color, light *LightRenderData, lightNormal, surfaceNormal *Point, material *Material) *colors.Color {
	dot := lightNormal.Dot(surfaceNormal)
	if dot <= 0.0 {
		return c
	}

	// Apply diffuse phong shading
	c = c.AddChannels(light.ColorIntensity.Scale(dot))

	// Compute and apply specular phong shading
	reflectionNormal := surfaceNormal.Scale(dot * 2.0).Subtract(lightNormal)
	specularIntensity := math.Pow(0.5 + reflectionNormal.Dot(MakePointZ()), material.SpecularExponent) / 255.0
	specularColor := material.SpecularColor.Scale(specularIntensity * light.Intensity)
	return c.AddChannels(specularColor)
}
