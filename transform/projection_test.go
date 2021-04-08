package transform

import (
	"math"
	"testing"

	"github.com/reactivego/seen/float"
)

func TestProjectionPerspectiveMat4x4(t *testing.T) {
	var (
		width  float64 = 4000
		height float64 = 3000
		near   float64 = 100
		far    float64 = 200
		right          = 0.5 * width
		left           = -right
		top            = 0.5 * height
		bottom         = -top
	)
	p := Projection{right, top, near, far}
	m := p.PerspectiveMat4x4()

	// Check that zero value for x,y,z returns the correct values
	x, y, z, w := m.Transform(0, 0, 0, 1)
	if !float.EqualPairs(x, 0, y, 0, z, -400, w, 0) {
		t.Fail()
	}

	// Check that positive edges are equal to w value
	x, y, z, w = m.Transform(right, top, -near, 1)
	if !float.EqualPairs(x, w, y, w, z, -w) {
		t.Fail()
	}

	// check that negative edges are equal to negative w value
	x, y, z, w = m.Transform(left, bottom, -near, 1)
	if !float.EqualPairs(x, -w, y, -w, z, -w) {
		t.Fail()
	}

	// check that clipping works by nudging outside the edge
	x, y, _, w = m.Transform(left-1.0, bottom-1.0, -near, 1)
	if x >= -w {
		t.Fail()
	}
	if y >= -w {
		t.Fail()
	}
}

func TestPerspectiveFunc(t *testing.T) {
	fovy := 135.0 * (math.Pi / 180.0)
	p := Perspective(fovy, 4.0/3.0, 100.0, 200.0)
	t.Log(p)

	// TODO: Actually test something
}
