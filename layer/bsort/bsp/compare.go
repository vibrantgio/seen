package bsp

import (
	"github.com/vibrantgio/seen/float"
)

type PlaneComparison int

const (
	Coplanar = PlaneComparison(iota)
	Before
	Behind
	Splits
)

func Compare(l, r Plane) PlaneComparison {
	if parallel := float.Equal(l.Normal.Dot(r.Normal), 1.0); parallel {
		d := l.Normal.Dot(l.Barycenter)
		dr := r.Normal.Dot(r.Barycenter)
		if float.Equal(d, dr) {
			return Coplanar
		} else {
			if d > dr {
				return Before
			} else {
				return Behind
			}
		}
	}
	d := l.Normal.Dot(l.Barycenter)
	var c PlaneComparison
	for _, p := range r.Points {
		dp := l.Normal.Dot(p)
		if !float.Equal(dp, d) {
			if d > dp {
				if c != Before {
					if c == Coplanar {
						c = Before
					} else {
						return Splits
					}
				}
			} else {
				if c != Behind {
					if c == Coplanar {
						c = Behind
					} else {
						return Splits
					}
				}
			}
		}
	}
	return c
}
