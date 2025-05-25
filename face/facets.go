package face

import "github.com/vibrantgio/seen/point"

type Facets []Facet

// FacesWith creates faces from points and facets. A facet is a list of indexes
// into the list of points. Note that a point referenced in multiple facets will
// be inserted into multiple faces. Because of this points shared by multiple
// faces will be transformed multiple times instead of only once.
func (facets Facets) FacesWith(points point.Points) (faces Faces) {
	faces = make(Faces, len(facets))
	for i, facet := range facets {
		face := &faces[i]
		face.Id = UniqueId()
		face.Options = make(map[string]string)
		for _, j := range facet {
			face.Points = append(face.Points, points[j])
		}
	}
	return
}
