package seen

type Polyhedron struct {
	V []Vertex
	P []Polygon
}

func (p *Polyhedron) SetPolygons(vertices []Vertex, polygons[][]int) {
	p.V = vertices
	p.P = make([]Polygon,len(polygons))
	for polyindex,polygon := range polygons {
		p.P[polyindex].V = make([]*Vertex,len(polygon))
		for i,vertindex := range polygon {
			p.P[polyindex].V[i] = &vertices[vertindex]
		}
	}
}
