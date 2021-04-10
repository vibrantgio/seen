package colors

import (
	"testing"

	"github.com/reactivego/seen/float"
)

func TestColorScale(t *testing.T) {

	c1 := White
	c2 := White.Scale(0.5)

	r1, g1, b1 := c1.R/2.0, c1.G/2.0, c1.B/2.0
	r2, g2, b2 := c2.R, c2.G, c2.B

	if !float.EqualPairs(r1, r2, g1, g2, b1, b2) {
		t.Logf("Exp {%v,%v,%v}", r1, g1, b1)
		t.Logf("Got {%v,%v,%v}", r2, g2, b2)
		t.Fail()
	}
}

func TestColorStringParsing(t *testing.T) {

	c, err := ColorWithString("#0F0F0F")
	if err != nil || !c.Equal(Color{15.0 / 255.0, 15.0 / 255.0, 15.0 / 255.0, 1}) {
		t.Log(c.R, c.G, c.B, c.A)
		t.Fail()
	}

	c, err = ColorWithString("#FFFFFF")
	if err != nil || !c.Equal(White) {
		t.Log(c.R, c.G, c.B, c.A)
		t.Fail()
	}
}

func TestColorStringGeneration(t *testing.T) {

	c := White
	if c.Hex() != "#FFFFFF" {
		t.Log(c.Hex())
		t.Fail()
	}

	c = Grey
	if c.Hex() != "#7F7F7F" {
		t.Log(c.Hex())
		t.Fail()
	}

	c = Black
	if c.Hex() != "#000000" {
		t.Log(c.Hex())
		t.Fail()
	}
}

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
