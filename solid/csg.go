package solid

import "slices"

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
	a, b := &BSP{}, &BSP{}
	a.AddPolygons(s.Clone())
	b.AddPolygons(other.Clone())
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
	a, b := &BSP{}, &BSP{}
	a.AddPolygons(s.Clone())
	b.AddPolygons(other.Clone())
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
	a, b := &BSP{}, &BSP{}
	a.AddPolygons(s.Clone())
	b.AddPolygons(other.Clone())
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
