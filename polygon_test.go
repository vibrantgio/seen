package seen

import (
	"testing"
	"xpt.nl/float"
)

func TestPolygonVertPointers(t *testing.T) {
	verts := [...]Vertex{ {1,2,3}, {4,5,6} }
	pvert := &verts[1]
	if !float.Equal(pvert.X, verts[1].X) {
		t.Fail()
	}
	pvert.X = 123	
	if !float.Equal(pvert.X, verts[1].X) {
		t.Fail()
	}
}

func TestPolygonCreation(t *testing.T) {
	p := &Polygon{}
	p.V = []*Vertex{ &Vertex{1,2,3}, &Vertex{4,5,6} }
}

func TestPolygonIcosahedron(t * testing.T) {
	// From seenjs.io
	ICOSAHEDRON_POLYGONS := [...][3]int{
		{0, 4, 1},
		{0, 9, 4},
		{9, 5, 4},
		{4, 5, 8},
		{4, 8, 1},
		{8, 10, 1},
		{8, 3, 10},
		{5, 3, 8},
		{5, 2, 3},
		{2, 7, 3},
		{7, 10, 3},
		{7, 6, 10},
		{7, 11, 6},
		{11, 0, 6},
		{0, 1, 6},
		{6, 1, 10},
		{9, 0, 11},
		{9, 11, 2},
		{9, 2, 5},
		{7, 2, 11},
	}
	const ICOS_X = 0.525731112119133606
	const ICOS_Z = 0.850650808352039932
	ICOSAHEDRON_VERTICES := [...]Vertex {
		{-ICOS_X, 0.0,     -ICOS_Z},
		{ICOS_X,  0.0,     -ICOS_Z},
		{-ICOS_X, 0.0,     ICOS_Z},
		{ICOS_X,  0.0,     ICOS_Z},
		{0.0,     ICOS_Z,  -ICOS_X},
		{0.0,     ICOS_Z,  ICOS_X},
		{0.0,     -ICOS_Z, -ICOS_X},
		{0.0,     -ICOS_Z, ICOS_X},
		{ICOS_Z,  ICOS_X,  0.0},
		{-ICOS_Z, ICOS_X,  0.0},
		{ICOS_Z,  -ICOS_X, 0.0},
		{-ICOS_Z, -ICOS_X, 0.0},
	}

	var icosahedron []Polygon
	for i := range ICOSAHEDRON_POLYGONS {
		p := Polygon{}
		for _,j := range ICOSAHEDRON_POLYGONS[i] {
			p.V = append(p.V, &ICOSAHEDRON_VERTICES[j])
		}
		icosahedron = append(icosahedron, p)
	}

	if len(icosahedron) != len(ICOSAHEDRON_POLYGONS) {
		t.Fail()
	}

	// Check that modifying a vertex in the global array propagates changes to the polygon.

	v := icosahedron[0].V[0]
	if !float.Equal(v.X,-ICOS_X) || !float.Equal(v.Y,0.0) || !float.Equal(v.Z,-ICOS_Z) {
		t.Fail()
	}
	w := &ICOSAHEDRON_VERTICES[0]
	w.X, w.Y, w.Z = 1,2,3
	if !float.Equal(v.X,1) || !float.Equal(v.Y,2) || !float.Equal(v.Z,3) {
		t.Fail()
	}


	// Check that the first vertex of every polygon has the correct values for X,Y,Z
	// for _,p := range icosahedron {
	// 	v := p.V[0]
	// 	if !float.Equal(v.X,-ICOS_X) || !float.Equal(v.Y,0.0) || !float.Equal(v.Z,-ICOS_Z) {
	// 		t.Fail()
	// 	}
	// }
	// // modify an original vertex
	// ICOSAHEDRON_VERTICES[0].X = 1
	// ICOSAHEDRON_VERTICES[0].Y = 2
	// ICOSAHEDRON_VERTICES[0].Z = 3
	// // check that all polygons see that change
	// for _,p := range icosahedron {
	// 	v := p.V[0]
	// 	if !float.Equal(v.X,1) || !float.Equal(v.Y,2) || !float.Equal(v.Z,3) {
	// 		t.Fail()
	// 	}
	// }
}
