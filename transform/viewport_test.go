package transform

import (
	"testing"
	"github.com/reactivego/seen/float"
)

func TestViewport(t *testing.T) {
	v := &Viewport{0, 0, 400, 300, 0, 1}

	x, y, z := v.Convert(-1, 1, -1)
	t.Log(x, y, z)
	if !float.Equal(x, 0) {
		t.Fail()
	}
	if !float.Equal(y, 0) {
		t.Fail()
	}
	if !float.Equal(z, 0) {
		t.Fail()
	}
	x, y, z = v.Convert(1, -1, 1)
	t.Log(x, y, z)
	if !float.Equal(x, 400) {
		t.Fail()
	}
	if !float.Equal(y, 300) {
		t.Fail()
	}
	if !float.Equal(z, 1) {
		t.Fail()
	}

	// Map project z with value 1 to far plane 1000
	v.F = 1000
	x, y, z = v.Convert(1, 1, 1)
	t.Log(x, y, z)
	if !float.Equal(x, 400) {
		t.Fail()
	}
	if !float.Equal(y, 0) {
		t.Fail()
	}
	if !float.Equal(z, 1000) {
		t.Fail()
	}
	_, _, z = v.Convert(1, 1, 0)
	t.Log(z)
	if !float.Equal(z, 500) {
		t.Fail()
	}

}
