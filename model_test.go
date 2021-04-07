package seen

import (
	"math"
	"testing"

	"github.com/reactivego/seen/transform"
)

var mock_ShapeCount int

func mock_ModelShapeFunc(shape *Shape, lights []LightShaderData, transform Matrix) {
	mock_ShapeCount++
}

func mock_MakeRectangle() *Shape {
	points := []Point{
		{0, 0, 0},
		{0, 0.5, 0},
		{0.5, 0, 0},
		{0.5, 0.5, 0},
	}
	coords := [][]int{
		{0, 1, 2},
		{2, 1, 3},
	}
	s := ShapeWith("rectangle", MakeSurfaces(points, coords))
	return &s
}

func mock_MakeText(message string) *Shape {
	points := []Point{
		{0, 0, 0},
		{0, 0.5, 0},
		{0.5, 0, 0},
	}
	s := MakeSurface(points)
	s.Options["text"] = message
	s.Init()
	t := ShapeWith("text", []Surface{*s})
	return &t
}

func TestModelAdding(t *testing.T) {
	m := MakeModel()
	s := mock_MakeRectangle()
	tx := mock_MakeText("Hello, World!")
	m2 := MakeModel()

	// Rotate around y axis (rhs coord system with +y pointing up,
	// +x pointing right and +z pointing out of the screen)
	r := transform.QuatAxisAngle(0, 1, 0, math.Pi/4.0)
	m2.SetRotation(r)
	m2.Add(s)
	m2.Add(tx)
	m.Add(s)
	m.Add(m2)

	m.EachRenderable(mock_ModelShapeFunc)

	if mock_ShapeCount != 3 {
		t.Fail()
	}
}
