package color

// SourceOption is the interface type for passing in color options to a source
// creation function.
type SourceOption interface {
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
