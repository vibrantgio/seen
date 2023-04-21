package color

import "math/rand"

// Source is an interface to a color source. It is used for generating
// a sequence of random colors that are slightly different.
type Source interface {
	Read() Color
}

// Option is the interface type for passing in color options for RandomSource2
type Option interface {
	Value() float64
}

// Drift default is 0.03
type Drift float64

func (v Drift) Value() float64 { return float64(v) }

// Hue default is a random value in the range [0-1]
type Hue float64

func (v Hue) Value() float64 { return float64(v) }

// Sat default is 0.5
type Sat float64

func (v Sat) Value() float64 { return float64(v) }

// Lit default is 0.4
type Lit float64

func (v Lit) Value() float64 { return float64(v) }

// Opacity default is 1.0
type Opacity float64

func (v Opacity) Value() float64 { return float64(v) }

// RandomSource generates a new random color every time it is called.
type RandomSource struct {
	sat, lit, opacity float64
}

// DefaultRandomSource generates a random hue every time a Color is read.
// It uses the default values Sat: 0.5, Lit: 0.4 and Opacity: 1.0 as parameters
// for the `ColorHsl()“ function.
func DefaultRandomSource() *RandomSource {
	return &RandomSource{
		sat:     0.5,
		lit:     0.4,
		opacity: 1.0,
	}
}

// RandomSourceWith is used to intialize the instance of a RandomColorSource with options.
func RandomSourceWith(options ...Option) Source {
	c := DefaultRandomSource()
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
	return ColorHsl(rand.Float64(), c.sat, c.lit, c.opacity)
}

// RandomSource2 is a source of random colors. It implements the Source interface.
type RandomSource2 struct {
	drift, hue, sat, lit, opacity float64
}

// DefaultRandomSource2 generates a random hue then randomly drifts the hue every
// time a Color is read. It uses the default values `Drift(0.03)`, `Sat(0.5)` and
// `Lit(0.4)` as input parameters for the hue drift algorithm. To start with a
// specific Hue, use the option e.g. `Hue(0.5)`.
func DefaultRandomSource2() *RandomSource2 {
	return &RandomSource2{
		drift:   0.03,
		hue:     rand.Float64(),
		sat:     0.5,
		lit:     0.4,
		opacity: 1.0,
	}
}

// RandomSource2With is used to intialize the instance of a RandomColorSource2 with options.
func RandomSource2With(options ...Option) Source {
	c := DefaultRandomSource2()
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

func (c *RandomSource2) Read() Color {
	c.hue += (rand.Float64() - 0.5) * c.drift
	for c.hue < 0 {
		c.hue += 1
	}
	for c.hue > 1 {
		c.hue -= 1
	}
	return ColorHsl(c.hue, c.sat, c.lit, c.opacity)
}
