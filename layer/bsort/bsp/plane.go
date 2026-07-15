package bsp

import (
	"github.com/vibrantgio/seen/float"
	"github.com/vibrantgio/seen/matrix"
	"github.com/vibrantgio/seen/point"
)

// Plane is the plane going through the Points of a face in the WorldSpace
// coordinate system.
type Plane struct {
	// Id of the face of the node this plane is associated with.
	Id int

	// Points contains the Face points in WorldSpace coordinates.
	Points point.Points

	// Barycenter contains the Barycenter of the Points (in WorldSpace
	// coordinates).
	Barycenter point.Point

	// Normal is the normal vector that is perpendicular to the plane going
	// through Points (in WorldSpace coordinates).
	Normal point.Point

	// Piece marks a plane produced by Split: a fragment of the face
	// identified by Id rather than the whole face. The renderer projects a
	// piece's own Points instead of the cached whole-face coordinates.
	Piece bool

	// NoSplit marks a plane whose polygon must never be cut (e.g. a text
	// face, which paints its whole string from its points — each piece
	// would repeat the text). Process routes a straddling NoSplit plane
	// wholesale to its barycenter's side instead of splitting it.
	NoSplit bool
}

func PlaneWith(id int, points point.Points, model matrix.Matrix) Plane {
	plane := Plane{Id: id, Points: make(point.Points, len(points))}
	plane.Barycenter = points.MulB(model, plane.Points)
	plane.Normal = plane.Points.Normal().Normalize()
	return plane
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
