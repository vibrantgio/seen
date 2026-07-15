package bsp

import (
	"math"

	"github.com/vibrantgio/seen/float"
)

// SideEpsilon is the absolute world-space distance within which a point is
// considered to lie on a partition plane. Plane normals are unit length, so
// Normal.Dot(p)-d is a true distance and an absolute tolerance keeps the
// classification independent of a model's distance from the origin (a
// relative epsilon's dead zone grows with coordinate magnitude). Measured on
// mocap scenes, genuine interpenetrations are >= 1e-3 world units deep while
// numerical jitter stays below 1e-6, so 1e-5 separates the two cleanly.
const SideEpsilon = 1e-5

// side classifies a signed distance from a plane: -1 below the plane (on the
// negative-normal side), +1 above it, 0 on the plane within SideEpsilon.
func side(distance float64) int {
	switch {
	case distance < -SideEpsilon:
		return -1
	case distance > +SideEpsilon:
		return +1
	case math.IsNaN(distance):
		// A NaN distance (degenerate plane normal) must not read as "on
		// the plane": grouping unrelated faces as coplanar with a garbage
		// partition would collapse them into one unsorted node. Pick a
		// side deterministically instead.
		return +1
	}
	return 0
}

type PlaneComparison int

const (
	Coplanar = PlaneComparison(iota)
	Before
	Behind
	Splits
)

// Compare classifies plane r against plane l: Coplanar when r lies in l's
// plane, Before when r is entirely on l's negative-normal side, Behind when
// entirely on the positive-normal side, and Splits when r has vertices on
// both sides (l's plane cuts r's polygon).
func Compare(l, r Plane) PlaneComparison {
	if parallel := float.Equal(l.Normal.Dot(r.Normal), 1.0); parallel {
		d := l.Normal.Dot(l.Barycenter)
		dr := r.Normal.Dot(r.Barycenter)
		switch side(d - dr) {
		case 0:
			return Coplanar
		case +1:
			return Before
		default:
			return Behind
		}
	}
	d := l.Normal.Dot(l.Barycenter)
	c := Coplanar
	for _, p := range r.Points {
		switch side(l.Normal.Dot(p) - d) {
		case -1:
			if c == Behind {
				return Splits
			}
			c = Before
		case +1:
			if c == Before {
				return Splits
			}
			c = Behind
		}
	}
	return c
}
