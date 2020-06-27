package colors

import (
	"errors"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"github.com/reactivego/seen/float"
)

// ColorReader is an interface used to read colors from a source.
// Used for reading of randomized colors for coloring surfaces.
type ColorReader interface {
	ReadColor() *Color
}

// Color objects store RGB and Alpha values with components in range [0..1]
type Color struct {
	R, G, B, A float64
}

var (
	// White is a shortcut for the white color
	White = &Color{1.0, 1.0, 1.0, 1.0}

	// Grey is a shortcut for the grey color.
	Grey = &Color{0.5, 0.5, 0.5, 1.0}

	// Black is a shortcut for the black color.
	Black = &Color{0.0, 0.0, 0.0, 1.0}
)

// MakeColorHsl creates a new `Color` using the supplied hue, saturation,
// and lightness (HSL) values.
// Each value must be in the range [0.0, 1.0].
func MakeColorHsl(h, s, l, a float64) *Color {
	// When saturation is 0, the color is "achromatic" or "grayscale".
	if s == 0 {
		return &Color{l, l, l, a}
	}
	var q float64
	if l < 0.5 {
		q = l * (1.0 + s)
	} else {
		q = l + s - l*s
	}
	p := 2.0*l - q

	r := hue2rgb(p, q, h+1.0/3.0)
	g := hue2rgb(p, q, h)
	b := hue2rgb(p, q, h-1.0/3.0)

	return &Color{r, g, b, a}
}

// Helper function to convert hue to rgb
func hue2rgb(p, q, t float64) float64 {
	switch {
	case t < 0.0:
		t += 1.0
	case t > 1.0:
		t -= 1.0
	}
	switch {
	case t < 1.0/6.0:
		return p + (q-p)*6.0*t
	case t < 1.0/2.0:
		return q
	case t < 2.0/3.0:
		return p + (q-p)*(2.0/3.0-t)*6.0
	default:
		return p
	}
}

// MakeColorWithString creates a color based on a string pattern #rrggbb or #rrggbbaa
// where the color components are hexadecimal values between 0 and 0xff. e.g. #ffffffff
// is white fully opaque
func MakeColorWithString(s string) (c *Color, err error) {
	if !strings.HasPrefix(s, "#") {
		return nil, errors.New("Parse Error: expected # as first character for color reference")
	}
	if len(s) != 9 && len(s) != 7 {
		return nil, errors.New("Parse Error: color reference does not match pattern #rrggbbaa or #rrggbb")
	}
	if len(s) != 9 {
		s += "FF"
	}
	h, err := strconv.ParseUint(s[1:], 16, 32)
	if err != nil {
		return nil, err
	}
	c = &Color{}
	c.R = float64(h&0xFF000000) / 0xFF000000
	c.G = float64(h&0x00FF0000) / 0x00FF0000
	c.B = float64(h&0x0000FF00) / 0x0000FF00
	c.A = float64(h&0x000000FF) / 0x000000FF
	return c, nil
}

// Hex returns a color string according to the following pattern #rrggbb
// Where the components are hexadecimal values between 0 and 0xff.
// The alpha component is left out of the string.
func (c Color) Hex() string {
	h := uint64(c.R*0xFF0000)&0xFF0000 +
		uint64(c.G*0x00FF00)&0x00FF00 +
		uint64(c.B*0x0000FF)&0x0000FF
	s := strings.ToUpper(strconv.FormatUint(h, 16))
	return "#000000"[:7-len(s)] + s
}

// Equal returns true when the colors have both equal color components as well as equal alpha value.
func (l *Color) Equal(r *Color) bool {
	return float.EqualPairs(l.R, r.R, l.G, r.G, l.B, r.B, l.A, r.A)
}

// Scale returns a Color with the rgb channels scaled by the supplied scalar value.
func (c Color) Scale(s float64) *Color {
	c.R *= s
	c.G *= s
	c.B *= s
	return &c
}

// AddChannels adds the channels of the current Color with each respective
// channel from the supplied Color object.
func (l Color) AddChannels(r *Color) *Color {
	l.R += r.R
	l.G += r.G
	l.B += r.B
	return &l
}

// MultiplyChannels multiplies the channels of the current Color with each respective
// channel from the supplied Color object.
func (l Color) MultiplyChannels(r *Color) *Color {
	l.R *= r.R
	l.G *= r.G
	l.B *= r.B
	return &l
}

// MinChannels takes the minimum between each channel of the current Color and the supplied Color object.
func (l Color) MinChannels(r *Color) *Color {
	l.R = math.Min(r.R, l.R)
	l.G = math.Min(r.G, l.G)
	l.B = math.Min(r.B, l.B)
	return &l
}

// Clamp clamps each rgb channel to the supplied minimum and maximum scalar values.
func (l Color) Clamp(min, max float64) *Color {
	l.R = math.Min(max, math.Max(min, l.R))
	l.G = math.Min(max, math.Max(min, l.G))
	l.B = math.Min(max, math.Max(min, l.B))
	return &l
}

// init will seed the default random generator with the current time.
// This is needed by the RandomColorReader2 struct below.
func init() {
	rand.Seed(time.Now().UnixNano())
}

// ColorOption is the interface type for passing in color options for RandomSurfaces2
type ColorOption interface{}

// ColorDrift default is 0.03
type ColorDrift float64

//ColorSat default is 0.5
type ColorSat float64

// ColorLit default is 0.4
type ColorLit float64

// RandomColorReader2 is a source of random colors. It implements the ColorReader interface.
type RandomColorReader2 struct {
	drift, sat, lit, hue float64
}

// MakeRandomColorReader2 generates a random hue then randomly drifts the hue every time a Color is read.
func MakeRandomColorReader2(options ...ColorOption) ColorReader {
	r := &RandomColorReader2{}
	r.Init()
	return r
}

// Init is used to intialize the instance of a RandomColorSource2
func (c *RandomColorReader2) Init(options ...ColorOption) {
	c.drift = 0.03
	c.sat = 0.5
	c.lit = 0.4
	for _, opt := range options {
		switch o := opt.(type) {
		case ColorDrift:
			c.drift = float64(o)
		case ColorSat:
			c.sat = float64(o)
		case ColorLit:
			c.lit = float64(o)
		}
	}
	c.hue = rand.Float64()
}

func (c *RandomColorReader2) ReadColor() *Color {
	c.hue += (rand.Float64() - 0.5) * c.drift
	for c.hue < 0 {
		c.hue += 1
	}
	for c.hue > 1 {
		c.hue -= 1
	}
	return MakeColorHsl(c.hue, c.sat, c.lit, 1.0)
}
