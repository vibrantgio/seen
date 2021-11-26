package bsp

import (
	"fmt"
	"os"

	"github.com/reactivego/seen"
	"github.com/reactivego/seen/float"
)

type PlaneComparison int

const (
	Coplanar = PlaneComparison(iota)
	Before
	Behind
	Splits
)

func Compare(l, r seen.Plane) PlaneComparison {
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

type Builder struct {
	TransformStack
	Planes []seen.Plane
}

func (v *Builder) VisitSurface(s *seen.Surface) {
	p := seen.Plane{Surface: s, Points: make(seen.Points, len(s.Points))}
	p.Barycenter = s.Points.Mul(v.Transform, p.Points)
	p.Normal = p.Points.Normal().Normalize()
	v.Planes = append(v.Planes, p)
}

func process(plane []seen.Plane, i int, recursion int, report func(...interface{})) *BSP {
	bsp := BSP{Plane: []seen.Plane{plane[i]}}
	planei := bsp.Plane[0]
	var before, behind []seen.Plane
	for j, planej := range plane {
		if j == i {
			continue
		}
		switch Compare(planei, planej) {
		case Coplanar:
			bsp.Plane = append(bsp.Plane, planej)
		case Before:
			before = append(before, planej)
		case Behind:
			behind = append(behind, planej)
		case Splits:
			if Compare(planej, planei) != Splits {
				// Situation: plane[i] splits plane[j] but not vice versa.
				if recursion < 16 {
					// Try  to partition space with plane[j] instead of plane[i]
					return process(plane, j, recursion+1, report)
				} else {
					// Situation: Were are probably looping....
					// TBD: use plane[i] to split plane[j]
					behind = append(behind, planej)
					if report != nil {
						report("split loop", i, j)
					}
				}
			} else {
				// Situation: planes[i] and planes[j] split each other.
				// TBD: use plane[i] to split plane[j]
				behind = append(behind, planej)
				if report != nil {
					report("split conflict", i, j)
					planej.Surface.FillMaterial, _ = seen.MaterialWith("#ff0000")
				}
			}
		}
	}
	if len(before) > 0 {
		bsp.Front = process(before, len(before)/2, 0, report)
	}
	if len(behind) > 0 {
		bsp.Back = process(behind, len(behind)/2, 0, report)
	}
	return &bsp
}

func (v *Builder) Build() *BSP {
	if len(v.Planes) == 0 {
		return nil
	}

	const NoReporting = true
	if NoReporting {
		return process(v.Planes, len(v.Planes)/2, 0, nil)
	}

	fmt.Fprintln(os.Stderr, "Building BSP...")
	bsp := process(v.Planes, len(v.Planes)/2, 0, func(p ...interface{}) {
		fmt.Fprintln(os.Stderr, p...)
	})
	fmt.Fprintln(os.Stderr, "BSP Done")
	return bsp
}
