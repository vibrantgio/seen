package bsp

import (
	"github.com/vibrantgio/seen/point"
)

// Split cuts target's polygon by the receiver's plane and returns the two
// pieces: front on the receiver's negative-normal side (the side the Front
// subtree holds, see Process) and back on the positive-normal side. Both
// pieces keep target's Id and Normal — they lie in the same plane and shade
// from the same face — and are marked as Piece so the renderer projects
// their own points instead of the cached whole-face coordinates.
//
// ok is false when target does not actually straddle the plane within
// SideEpsilon or a piece would degenerate to fewer than three vertices;
// callers should then keep the whole target on a single side.
//
// The cut walks the polygon ring once, so a concave polygon that crosses
// the plane more than twice comes back as two self-touching rings rather
// than several disjoint pieces. Faces in practice (pipe quads, sphere
// triangles) are convex, where a single cut yields exactly two pieces.
func (l Plane) Split(target Plane) (front, back Plane, ok bool) {
	d := l.Normal.Dot(l.Barycenter)

	n := len(target.Points)
	distance := make([]float64, n)
	sides := make([]int, n)
	hasNeg, hasPos := false, false
	for i, p := range target.Points {
		distance[i] = l.Normal.Dot(p) - d
		sides[i] = side(distance[i])
		hasNeg = hasNeg || sides[i] < 0
		hasPos = hasPos || sides[i] > 0
	}
	if !hasNeg || !hasPos {
		return front, back, false
	}

	var neg, pos point.Points
	for i := 0; i < n; i++ {
		switch {
		case sides[i] < 0:
			neg = append(neg, target.Points[i])
		case sides[i] > 0:
			pos = append(pos, target.Points[i])
		default:
			// On the plane: the cut runs through this vertex, so both
			// pieces share it.
			neg = append(neg, target.Points[i])
			pos = append(pos, target.Points[i])
		}
		// A strict sign change along the edge to the next vertex crosses
		// the plane: insert the intersection point into both rings.
		if j := (i + 1) % n; sides[i]*sides[j] < 0 {
			t := distance[i] / (distance[i] - distance[j])
			x := target.Points[i].Plus(target.Points[j].Minus(target.Points[i]).Times(t))
			neg = append(neg, x)
			pos = append(pos, x)
		}
	}
	if len(neg) < 3 || len(pos) < 3 {
		return front, back, false
	}

	front = Plane{Id: target.Id, Points: neg, Barycenter: neg.Barycenter(), Normal: target.Normal, Piece: true}
	back = Plane{Id: target.Id, Points: pos, Barycenter: pos.Barycenter(), Normal: target.Normal, Piece: true}
	return front, back, true
}
