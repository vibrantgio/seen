package seen

import (
	"math"
)

type Bounds struct {
	Min, Max *Point
}

func MakeBounds(points []Point) *Bounds {
	b := &Bounds{}
	for i := range points {
		b.AddAssign(&points[i])
	}
	return b
}

func (b Bounds) Add(p *Point) *Bounds {
	b.AddAssign(p)
	return &b
}

func (b *Bounds) AddAssign(p *Point) {
	if b.Min == nil {
		min := *p // force copy
		b.Min = &min
	} else {
		b.Min.X = math.Min(b.Min.X, p.X)
		b.Min.Y = math.Min(b.Min.Y, p.X)
		b.Min.Z = math.Min(b.Min.Z, p.X)
	}

	if b.Max == nil {
		max := *p // force copy
		b.Max = &max
	} else {
		b.Max.X = math.Max(b.Max.X, p.X)
		b.Max.Y = math.Max(b.Max.Y, p.Y)
		b.Max.Z = math.Max(b.Max.Z, p.Z)
	}
}
