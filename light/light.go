package light

import (
	"github.com/vibrantgio/seen/color"
	"github.com/vibrantgio/seen/intensity"
	"github.com/vibrantgio/seen/matrix"
	"github.com/vibrantgio/seen/point"
	"github.com/vibrantgio/seen/transform"
)

// Light object holds the attributes and transformation of a light source.
type Light struct {
	transform.Transform
	Kind  Kind
	Point point.Point

	// Color is the color of the light.
	Color color.Color

	// Intensity should be a value between 0.0 and 1.0 that determines the
	// ammount of light contributed by this light. An intensity of 0.0
	// effectively turns the light off while a value of 1.0 will add the
	// value of the Color field to the face being lit.
	Intensity intensity.Intensity

	Normal  point.Point
	Enabled bool
}

func Of(kind Kind) (l Light) {
	l.Transform = transform.Default
	l.Kind = kind
	l.Point = point.Point{}
	l.Color = color.White
	l.Intensity = intensity.Default
	l.Normal = point.Point{X: 1, Y: -1, Z: -1}.Normalize()
	l.Enabled = true
	return
}

func (l Light) IsEnabled() bool {
	return l.Enabled
}

// ShaderData pre-computes values necessary for shading.
func (l Light) ShaderData(model matrix.Matrix) (lsd ShaderData) {
	lsd.Kind = l.Kind
	lsd.Point = l.Point.Mul(model)
	lsd.Color = l.Color.Scale(float64(l.Intensity))
	lsd.Intensity = l.Intensity
	origin := point.Point{}.Mul(model)
	lsd.Normal = l.Normal.Mul(model).Minus(origin).Normalize()
	return lsd
}
