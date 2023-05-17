package seen

import "github.com/reactivego/seen/color"

// Material objects hold the attributes that desribe the color and finish of a surface.
type Material struct {
	// The base color of the material.
	Color color.Color

	// Metallic property determines how the specular highlights are
	// calculated. Normally, specular highlights are the color of the light
	// source. If Metallic is true, specular highlight colors are determined
	// from the SpecularColor property.
	Metallic bool

	// The color used for specular highlights when `metallic` is true.
	SpecularColor color.Color

	// SpecularExponent determines how "shiny" the material is. A low
	// exponent will create a low-intesity, diffuse specular shine. A high
	// exponent will create an intense, point-like specular shine.
	SpecularExponent float64

	// Shader object may be supplied to override the shader used for this
	// material. For example, if you want to apply a flat color to text or
	// other shapes, set this value to FlatShader.
	Shader Shader
}

// NewMaterialWith makes a material based on the given source parameter.
// The source can be another Material, Color, or string containing a hex color representation.
func NewMaterialWith(source interface{}) (m *Material, err error) {
	err = nil
	switch s := source.(type) {
	case Material:
		m = &s
	case *Material:
		mc := *s
		m = &mc
	case color.Color:
		m = &Material{
			Color:            s,
			SpecularColor:    color.White,
			SpecularExponent: 15.0,
		}
	case string:
		c, err := color.ColorWithString(s)
		if err == nil {
			m = &Material{
				Color:            c,
				SpecularColor:    color.White,
				SpecularExponent: 15.0,
			}
		}
	default:
		m = &Material{
			Color:            color.Grey,
			SpecularColor:    color.White,
			SpecularExponent: 15.0,
		}
	}
	return
}

// Render applies the shader's shading to this material, with the option to override
// the shader with the material's shader (if defined).
func (m *Material) Render(lights []LightShaderData, shader Shader, surface *SurfaceShaderData) color.Color {
	var fill color.Color
	if m.Shader != nil {
		fill = m.Shader.Shade(lights, surface, m)
	} else {
		fill = shader.Shade(lights, surface, m)
	}
	fill.A = m.Color.A
	return fill
}
