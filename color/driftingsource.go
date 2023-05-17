package color

import "math/rand"

// DriftingSource is a source of colors where each newly returned
// color has a slightly different Hue w.r.t. the previously generated
// color. The Drift (default 0.03) inidicates how much the hue is allowed
// to drift from the previous value. The Sat (default 0.5) and Lit (default 0.4)
// indicate the saturation and lightness of the generated colors. The Opacity
// (default 1.0) indicates the opacity of the generated colors. The initial
// hue is randomly generated.
type DriftingSource struct {
	drift, hue, sat, lit, opacity float64
}

var _ Source = (*DriftingSource)(nil)

// NewDriftingSource creates a color source that returns a sequence of colors where
// the hue of each color is slightly different from the previous color.
// The color source's options are by default:
//
//	Drift(0.03)
//	Hue(rand.Float64())
//	Sat(0.5)
//	Lit(0.4)
//	Opacity(1.0)
func NewDriftingSource() *DriftingSource {
	return &DriftingSource{
		drift:   0.03,
		hue:     rand.Float64(),
		sat:     0.5,
		lit:     0.4,
		opacity: 1.0,
	}
}

// NewDriftingSourceWith createsa a color source that returns a sequence of colors where
// the hue of each color is slightly different from the previous color.
// The options of the color source can be set by passing in options to the function.
// The default options are as follows:
//
//	Drift(0.03)
//	Hue(rand.Float64())
//	Sat(0.5)
//	Lit(0.4)
//	Opacity(1.0)
func NewDriftingSourceWith(options ...SourceOption) *DriftingSource {
	c := NewDriftingSource()
	for _, opt := range options {
		switch o := opt.(type) {
		case Drift:
			c.drift = o.Value()
		case Hue:
			c.hue = o.Value()
		case Sat:
			c.sat = o.Value()
		case Lit:
			c.lit = o.Value()
		case Opacity:
			c.opacity = o.Value()
		}
	}
	return c
}

func (c *DriftingSource) Read() Color {
	c.hue += (rand.Float64() - 0.5) * c.drift
	for c.hue < 0 {
		c.hue += 1
	}
	for c.hue > 1 {
		c.hue -= 1
	}
	return ColorHSL(c.hue, c.sat, c.lit, c.opacity)
}
