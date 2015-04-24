package seen

type Polygon struct {
	// Pointers to vertices that form the outline of this polygon.
	// Vertices themselves reside in containing object so all vertices 
	// belonging to one object can be quickly transformed in one go.
	// Also because vertices are shared by polygons where their edges
	// touch, those vertices only need to transformed once.
	V []*Vertex
}
