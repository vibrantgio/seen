package layer

import (
	"github.com/vibrantgio/seen/color"
	"github.com/vibrantgio/seen/point"
)

// Fragment contains only the information needed to render the face on a canvas.
type Fragment struct {
	// Points of the fragment in screen coordinates
	Points point.Points
	// Fill color if assigned is result of shading by taking into account
	// material, lights and shading algorithm.
	Fill *color.Color
	// Stroke color if assigned is result of shading by taking into account
	// material, lights and shading algorithm.
	Stroke *color.Color
	// Options contains all the style options passed in on  the face.
	Options map[string]string
}
