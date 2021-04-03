package seen

import "github.com/reactivego/seen/colors"

// Light model object holds the attributes and transformation of a light source.
type Light struct {
	Object
	Kind  string
	Id    string
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

func MakeLight(kind string) *Light {
	l := &Light{}
	l.Init(kind)
	return l
}

func (l *Light) Init(kind string) {
	l.Object.Init()
	l.Kind = kind
	l.Id = UniqueId("l")
	l.Point = PointZero
	l.Color = colors.White
	l.Intensity = 0.5
	l.Normal = Point{1, -1, -1}
	l.Normal.Normalize()
	l.Enabled = true
}

// MakePointLight() makes a Light that emits light in all directions from a single point.
// The Point property determines the location of the point light. Note,
// though, that it may also be moved through the transformation of the light.
func MakePointLight() *Light {
	return MakeLight("point")
}

// MakeDirectionalLight() makes a light that emits light in parallel lines,
// not eminating from any single point. For these lights, only the Normal
// property is used to determine the direction of the light. This may also
// be transformed.
func MakeDirectionalLight() *Light {
	return MakeLight("directional")
}

// MakeAmbientLight() makes a light that emits a constant amount of light
// everywhere at once. Transformation of the light has no effect.
func MakeAmbientLight() *Light {
	return MakeLight("ambient")
}

// LightRenderData stores pre-computed values necessary for shading
// surfaces with the supplied Light
type LightRenderData struct {
	Light          *Light
	ColorIntensity colors.Color
	Kind           string
	Intensity      float64
	Point          Point
	Normal         Point
}

func MakeLightRenderData(light *Light, transform Matrix) *LightRenderData {
	l := &LightRenderData{}
	l.Init(light, transform)
	return l
}

func (l *LightRenderData) Init(light *Light, transform Matrix) {
	l.Light = light
	l.ColorIntensity = light.Color.Scale(light.Intensity)
	l.Kind = light.Kind
	l.Intensity = light.Intensity
	l.Point = transform.TransformPoint(light.Point)
	origin := transform.TransformPoint(PointZero)
	l.Normal = transform.TransformPoint(light.Normal).Subtract(origin).Normalize()
}
