package seen

import "github.com/reactivego/seen/float"

// Plane is the plane going through the Points of a
// surface in the WorldSpace coordinate system.
type Plane struct {
	// Surface points to the Surface this Plane represents.
	*Surface

	// Points contains the Surface points in WorldSpace coordinates.
	Points Points

	// Barycenter contains the Barycenter of the Points.
	Barycenter Point

	// Normal is the normal vector that is perpendicular
	// to the plane going through Points.
	Normal Point
}

// ParallelWith returns true when this plane is parallel to the given plane.
// This means the planes have the same normal, but are in parallel planes.
func (l Plane) ParallelWith(r Plane) bool {
	return float.Equal(l.Normal.Dot(r.Normal), 1.0)
}

// CoplanarWith returns true when this plane is in the same plane as the given plane.
func (l Plane) CoplanarWith(r Plane) bool {
	if parallel := float.Equal(l.Normal.Dot(r.Normal), 1.0); parallel {
		return float.Equal(l.Normal.Dot(l.Barycenter), r.Normal.Dot(r.Barycenter))
	}
	return false
}
