package shape

import (
	"github.com/vibrantgio/seen"
	"github.com/vibrantgio/seen/face"
	"github.com/vibrantgio/seen/point"
)

// Returns a shape that is an extrusion of the supplied points into the z axis.
func Extrude(points point.Points, offset point.Point) seen.Object {
	n := len(points)
	points = append(points, make(point.Points, n)...)
	for i := range n {
		points[n+i] = points[i].Plus(offset)
	}
	var facets face.Facets
	for i := range n {
		j := (i + 1) % n
		facets = append(facets, face.Facet{i, n + i, n + j, j})
	}
	var front face.Facet
	for i := range n {
		front = append(front, i)
	}
	var back face.Facet
	for i := n - 1; i >= 0; i-- {
		back = append(back, i)
	}
	return NewShape("extrusion", points, append(facets, front, back))
}
