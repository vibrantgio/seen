package face

import (
	"github.com/vibrantgio/seen/matrix"
	"github.com/vibrantgio/seen/point"
	"github.com/vibrantgio/seen/shader"
)

// Face is a defined as a planar object in 3D space. These paths don't
// necessarily need to be convex, but they should be non-degenerate.
type Face struct {
	// Points contain a list of vertices of the planar polygon that defines the
	// outline of the face.
	Points point.Points

	// Id holds a unique identifier for the face.
	// We store a unique Id for every face so we can look them up quickly
	// with the render face cache.
	Id int

	// ShowBackfaces when set to true will override backface culling, which is useful if your
	// material is transparent. See comment in Scene.
	ShowBackfaces bool

	// FillMaterial may be a Material object which defines the color and
	// finish of the object and are rendered using the scene's shader.
	// If not material is set a Material(C.gray) will be used.
	FillMaterial *shader.Material

	// StrokeMaterial may be a Material object that defines the color when
	// an object is stroked. By default no stroke material will be set.
	StrokeMaterial *shader.Material

	// Dirty flag can be set to force the Coordinates generated from the face to
	// be regenerated.
	Dirty bool

	// Options is a map of additional options that can be specified for a face.
	// The option with key "stroke-width" is passed in the style map parameter to
	// PathPainter.Stroke() call.
	// The keys "font" and "anchor" are passed in as keys "font" and "text-anchor" in
	// the style map parameter to TextPainter.FillText() call.
	Options map[string]string
}

// FaceWith takes a slice of points and uses it in constructing a face. The
// array backing this points slice is integrated into the face.
func FaceWith(points point.Points) Face {
	face := Face{}
	face.Id = UniqueId()
	face.Options = make(map[string]string)
	face.Points = points
	return face
}

func (face *Face) SetFill(value any) (err error) {
	face.FillMaterial, err = shader.NewMaterialWith(value)
	return
}

func (face *Face) SetStroke(value any) (err error) {
	face.StrokeMaterial, err = shader.NewMaterialWith(value)
	return
}

func (face *Face) Coordinates(model, projection, viewport matrix.Matrix) Coordinates {
	rendering := Coordinates{}
	rendering.Face = face
	rendering.Update(face.Points, model, projection, viewport)
	return rendering
}
