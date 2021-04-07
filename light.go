package seen

import "github.com/reactivego/seen/colors"

// Light model object holds the attributes and transformation of a light source.
type Light struct {
	Object
	Kind  string
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

func light(kind string) (l Light) {
	l.Object.Init()
	l.Kind = kind
	l.Point = PointZero
	l.Color = colors.White
	l.Intensity = 0.5
	l.Normal = Point{1, -1, -1}
	l.Normal.Normalize()
	l.Enabled = true
	return
}

// LightRenderData stores pre-computed values necessary for shading
// surfaces with the supplied Light
type LightRenderData struct {
	Kind      string
	Point     Point
	Color     colors.Color
	Intensity float64
	Normal    Point
}

func (l Light) RenderData(transform Matrix) (lrd LightRenderData) {
	lrd.Kind = l.Kind
	lrd.Point = transform.TransformPoint(l.Point)
	lrd.Color = l.Color.Scale(l.Intensity)
	lrd.Intensity = l.Intensity
	origin := transform.TransformPoint(PointZero)
	lrd.Normal = transform.TransformPoint(l.Normal).Subtract(origin).Normalize()
	return lrd
}

// PointLight is a Light that emits light in all directions from a single point.
// The Point property determines the location of the point light. Note,
// though, that it may also be moved through the transformation of the light.
var PointLight = light("point")

// DirectionalLight is a light that emits light in parallel lines,
// not eminating from any single point. For these lights, only the Normal
// property is used to determine the direction of the light. This may also
// be transformed.
var DirectionalLight = light("directional")

// AmbientLight is a light that emits a constant amount of light
// everywhere at once. Transformation of the light has no effect.
var AmbientLight = light("ambient")
