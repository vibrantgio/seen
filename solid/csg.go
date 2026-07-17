package solid

import (
	"math"
	"slices"
)

// CSG holds a binary space partition tree representing a 3D solid.
// Two solids can be combined using the `Union()`, `Subtract()` and
// `Intersect()` methods.
type CSG []Polygon

// Union returns a new CSG solid representing space in either this solid or
// in the solid `csg`. Neither this solid nor the solid `csg` are modified.
//
//	A.Union(B)
//
//	+-------+            +-------+
//	|       |            |       |
//	|   A   |            |       |
//	|    +--+----+   =   |       +----+
//	+----+--+    |       +----+       |
//	     |   B   |            |       |
//	     |       |            |       |
//	     +-------+            +-------+
func (s CSG) Union(other CSG) CSG {
	// Empty operands are handled algebraically: an empty BSP tree cannot
	// represent the difference between "empty solid" and "universe", so
	// inverting one inside the boolean pipeline is meaningless (and used to
	// panic on the nil root plane). Emptiness is judged AFTER dropping
	// degenerate polygons — a solid of zero-area debris is empty too.
	sp, op := s.cloneSpanning(), other.cloneSpanning()
	if len(sp) == 0 {
		return op
	}
	if len(op) == 0 {
		return sp
	}
	a, b := &BSP{}, &BSP{}
	a.AddPolygons(sp)
	b.AddPolygons(op)
	a.ClipTo(b)
	b.ClipTo(a)
	b.Invert()
	b.ClipTo(a)
	b.Invert()
	a.AddPolygons(b.AllPolygons())
	return CSG(a.AllPolygons())
}

// Subtract returns a new CSG solid representing space in this solid but not
// in the solid `csg`. Neither this solid nor the solid `csg` are modified.
//
//	A.Subtract(B)
//
//	+-------+            +-------+
//	|       |            |       |
//	|   A   |            |       |
//	|    +--+----+   =   |    +--+
//	+----+--+    |       +----+
//	     |   B   |
//	     |       |
//	     +-------+
func (s CSG) Subtract(other CSG) CSG {
	// Empty minus anything is empty; anything minus empty is unchanged.
	// (See Union for why empty operands must not enter the BSP pipeline.)
	sp, op := s.cloneSpanning(), other.cloneSpanning()
	if len(sp) == 0 {
		return CSG{}
	}
	if len(op) == 0 {
		return sp
	}
	a, b := &BSP{}, &BSP{}
	a.AddPolygons(sp)
	b.AddPolygons(op)
	a.Invert()
	a.ClipTo(b)
	b.ClipTo(a)
	b.Invert()
	b.ClipTo(a)
	b.Invert()
	a.AddPolygons(b.AllPolygons())
	a.Invert()
	return CSG(a.AllPolygons())
}

// Intersect returns a new CSG solid representing space both this solid and in
// the solid `csg`. Neither this solid nor the solid `csg` are modified.
//
//	A.intersect(B)
//
//	+-------+
//	|       |
//	|   A   |
//	|    +--+----+   =   +--+
//	+----+--+    |       +--+
//	     |   B   |
//	     |       |
//	     +-------+
func (s CSG) Intersect(other CSG) CSG {
	// Intersecting with emptiness is empty. (See Union for why empty
	// operands must not enter the BSP pipeline.)
	sp, op := s.cloneSpanning(), other.cloneSpanning()
	if len(sp) == 0 || len(op) == 0 {
		return CSG{}
	}
	a, b := &BSP{}, &BSP{}
	a.AddPolygons(sp)
	b.AddPolygons(op)
	a.Invert()
	b.ClipTo(a)
	b.Invert()
	a.ClipTo(b)
	b.ClipTo(a)
	a.AddPolygons(b.AllPolygons())
	a.Invert()
	return CSG(a.AllPolygons())
}

// Inverse returns a new CSG solid with solid and empty space switched.
// This solid is not modified.
func (s CSG) Inverse() CSG {
	polygons := s.Clone()
	for i := range polygons {
		polygons[i].Flip()
	}
	return CSG(polygons)
}

// Clone returns a new CSG solid that is a deep clone.
func (p CSG) Clone() CSG {
	p = slices.Clone(p)
	for i := range p {
		p[i].Vertices = slices.Clone(p[i].Vertices)
	}
	return p
}

// cloneSpanning deep-clones the polygons that span a real plane, dropping
// degenerate ones (collinear or duplicated points, recognisable by their NaN
// plane normal). Degenerate polygons have zero area, so removing them never
// changes the solid — but letting one into the boolean pipeline does: every
// dot product against a NaN partition plane is NaN, every comparison false,
// so it would collapse whole polygon sets into one unsorted node and make
// later clips keep everything.
func (p CSG) cloneSpanning() CSG {
	out := make(CSG, 0, len(p))
	for _, poly := range p {
		if math.IsNaN(poly.Plane.Normal.X) {
			continue
		}
		poly.Vertices = slices.Clone(poly.Vertices)
		out = append(out, poly)
	}
	return out
}
