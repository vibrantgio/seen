package shape

import (
	"slices"

	"github.com/vibrantgio/seen"
	"github.com/vibrantgio/seen/affine"
	"github.com/vibrantgio/seen/face"
)

// Text creates a Shape for the given text.
// The text parameter is set as face option "text" on the created shape.
// After the text is set the other face options that were passed in via faceOptions
// are assigned to the face options of the created shape.
//
//	font	e.g. "20px sans-serif" or "10px Roboto"
//	anchor	e.g. "middle"
func Text(text string, faceOptions map[string]string) seen.Object {
	f := face.FaceWith(slices.Clone(affine.ORTHONORMAL_BASIS[:]))
	f.Options["text"] = text
	for key, val := range faceOptions {
		f.Options[key] = val
	}
	return NewShapeWithFaces("text", face.Faces{f})
}
