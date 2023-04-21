package color

import "testing"

func TestDefaultRandomSource2(t *testing.T) {

	r := DefaultRandomSource2()
	c0 := r.Read()
	c1 := r.Read()
	c2 := r.Read()
	c3 := r.Read()
	if c0.Equal(c1) && c1.Equal(c2) && c2.Equal(c3) {
		t.Fail()
	}
}
