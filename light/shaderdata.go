package light

import (
	"github.com/vibrantgio/seen/color"
	"github.com/vibrantgio/seen/intensity"
	"github.com/vibrantgio/seen/point"
)

// ShaderData stores pre-computed values necessary for shading
// with a certain light.
type ShaderData struct {
	Kind      Kind
	Point     point.Point
	Color     color.Color
	Intensity intensity.Intensity
	Normal    point.Point
}
