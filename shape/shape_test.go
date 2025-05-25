package shape

import (
	"math"
	"testing"

	"github.com/vibrantgio/seen"
	"github.com/vibrantgio/seen/face"
	"github.com/vibrantgio/seen/light"
	"github.com/vibrantgio/seen/matrix"
	"github.com/vibrantgio/seen/point"
	"github.com/vibrantgio/seen/quaternion"
)

var mock_ShapeCount int

func mock_GroupShapeFunc(seen.Object, []light.ShaderData, matrix.Matrix) {
	mock_ShapeCount++
}

func mock_Rectangle() seen.Object {
	points := [...]point.Point{
		{0, 0, 0},
		{0, 0.5, 0},
		{0.5, 0, 0},
		{0.5, 0.5, 0},
	}
	facets := [...]face.Facet{
		{0, 1, 2},
		{2, 1, 3},
	}
	return NewShape("rectangle", points[:], facets[:])
}

func mock_Text(message string) seen.Object {
	points := [...]point.Point{
		{0, 0, 0},
		{0, 0.5, 0},
		{0.5, 0, 0},
	}
	f := face.FaceWith(points[:])
	f.Options["text"] = message
	return NewShapeWithFaces("text", face.Faces{f})
}

func TestGroupAdding(t *testing.T) {
	s := mock_Rectangle()
	tx := mock_Text("Hello, World!")

	// Rotate around y axis (rhs coord system with +y pointing up,
	// +x pointing right and +z pointing out of the screen)
	r := quaternion.AxisAngle(0, 1, 0, math.Pi/4.0)
	m2 := seen.NewGroup(s, tx)
	m2.SetRotation(r)
	m := seen.NewGroup(s, m2)

	m.EachRenderable(mock_GroupShapeFunc)

	if mock_ShapeCount != 3 {
		t.Fail()
	}
}
