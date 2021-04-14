package shape

import (
	"github.com/reactivego/seen"
	"github.com/reactivego/seen/affine"
)

// Text creates a Shape for the given text.
// The text parameter is set as surface option "text" on the created shape.
// After the text is set the other surface options that were passed in via surfaceOptions
// are assigned to the surface options of the created shape.
//	font	e.g. "20px sans-serif" or "10px Roboto"
//	anchor	e.g. "middle"
func Text(text string, surfaceOptions map[string]string) *seen.Shape {
	surface := seen.SurfaceWith(affine.ORTHONORMAL_BASIS)
	surface.Options["text"] = text
	for key, val := range surfaceOptions {
		surface.Options[key] = val
	}
	return &seen.Shape{
		Type:      "text",
		Transform: seen.DefaultTransform,
		Surfaces:  seen.Surfaces{*surface}}
}
