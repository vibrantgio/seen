package light

import (
	"github.com/vibrantgio/seen/color"
	"github.com/vibrantgio/seen/intensity"
	"github.com/vibrantgio/seen/point"
)

// DefaultLights are a set of lights to setup a standard Hollywood-style 3-part lighting
func DefaultLights() []*Light {
	// Key light
	key := DirectionalLight()
	key.Normal = point.Point{X: -1, Y: 1, Z: 1}.Normalize()
	key.Color = color.ColorHSL(0.1, 0.3, 0.7, 1.0)
	key.Intensity = intensity.Key

	// Back light
	back := DirectionalLight()
	back.Normal = point.Point{X: 1, Y: 1, Z: -1}.Normalize()
	back.Intensity = intensity.Back

	// Fill light
	fill := AmbientLight()
	fill.Intensity = intensity.Fill

	return []*Light{&key, &back, &fill}
}
