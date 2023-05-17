package color

import "math/rand"

// RandomSource generates a new random color every time it is called.
type RandomSource struct {
	sat, lit, opacity float64
}

var _ Source = (*RandomSource)(nil)

// NewRandomSource generates a random hue every time a Color is read.
// It uses the default values Sat: 0.5, Lit: 0.4 and Opacity: 1.0 as parameters
// for the `ColorHsl()“ function.
func NewRandomSource() *RandomSource {
	return &RandomSource{
		sat:     0.5,
		lit:     0.4,
		opacity: 1.0,
	}
}

// NewRandomSourceWith is used to intialize the instance of a RandomColorSource with options.
func NewRandomSourceWith(options ...SourceOption) *RandomSource {
	c := NewRandomSource()
	for _, opt := range options {
		switch o := opt.(type) {
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

func (c *RandomSource) Read() Color {
	return ColorHSL(rand.Float64(), c.sat, c.lit, c.opacity)
}
