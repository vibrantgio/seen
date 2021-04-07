package seen

import "github.com/reactivego/seen/colors"

type LightKind string

// Light model object holds the attributes and transformation of a light source.
type Light struct {
	Object
	Kind  LightKind
	Point Point

	// Color is the color of the light.
	Color colors.Color

	// Intensity should be a value between 0.0 and 1.0 that determines the
	// ammount of light contributed by this light. An intensity of 0.0
	// effectively turns the light off while a value of 1.0 will add the
	// value of the Color field to the surface being lit.
	Intensity float64

	Normal  Point
	Enabled bool
}

func LightWith(kind LightKind) (l Light) {
	l.Object = DefaultObject
	l.Kind = kind
	l.Point = PointZero
	l.Color = colors.White
	l.Intensity = 0.5 // 0.01
	l.Normal = Point{1, -1, -1}.Normalize()
	l.Enabled = true
	return
}

// LightShaderData stores pre-computed values necessary for shading
// with a certain light.
type LightShaderData struct {
	Kind      LightKind
	Point     Point
	Color     colors.Color
	Intensity float64
	Normal    Point
}

// ShaderData pre-computes values necessary for shading.
func (l Light) ShaderData(transform Matrix) (lsd LightShaderData) {
	lsd.Kind = l.Kind
	lsd.Point = transform.TransformPoint(l.Point)
	lsd.Color = l.Color.Scale(l.Intensity)
	lsd.Intensity = l.Intensity
	origin := transform.TransformPoint(PointZero)
	lsd.Normal = transform.TransformPoint(l.Normal).Subtract(origin).Normalize()
	return lsd
}

// PointLight is a Light that emits light in all directions from a single point.
// The Point property determines the location of the point light. Note,
// though, that it may also be moved through the transformation of the light.
func PointLight() Light {
	return LightWith(LightKind("point"))
}

// DirectionalLight is a light that emits light in parallel lines,
// not eminating from any single point. For these lights, only the Normal
// property is used to determine the direction of the light. This may also
// be transformed.
func DirectionalLight() Light {
	return LightWith(LightKind("directional"))
}

// AmbientLight is a light that emits a constant amount of light
// everywhere at once. Transformation of the light has no effect.
func AmbientLight() Light {
	return LightWith(LightKind("ambient"))
}

// DefaultLights are a set of lights to setup a standard Hollywood-style 3-part lighting
func DefaultLights() []Transformable {
	kl := DirectionalLight()
	kl.Normal = Point{-1, 1, 1}.Normalize()
	kl.Color = colors.ColorHsl(0.1, 0.3, 0.7, 1.0)
	kl.Intensity = 1.0 // 0.004

	// Back light
	bl := DirectionalLight()
	bl.Normal = Point{1, 1, -1}.Normalize()
	bl.Intensity = 0.765 // 0.003

	// Fill light
	fl := AmbientLight()
	fl.Intensity = 0.3825 // 0.0015

	return []Transformable{&kl, &bl, &fl}
}
