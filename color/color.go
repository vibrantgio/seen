package color

import (
	"image/color"
	"math"
	"strconv"
	"strings"

	"github.com/vibrantgio/seen/float"
)

// Color objects store RGB and Alpha values with components in range [0..1]
type Color struct {
	R, G, B, A float64
}

// ColorHSL creates a new `Color` using the supplied hue, saturation,
// and lightness (HSL) values.
// Each value must be in the range [0.0, 1.0].
func ColorHSL(h, s, l, a float64) Color {
	// When saturation is 0, the color is "achromatic" or "grayscale".
	if s == 0 {
		return Color{l, l, l, a}
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

	return Color{r, g, b, a}
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

// ColorWithString creates a color based on a string pattern #rgb, #rgba, #rrggbb
// or #rrggbbaa where the color components are hexadecimal values between 0 and 0xff.
// e.g. #ffffffff is white fully opaque.
func ColorWithString(s string) (c Color, err error) {
	if !strings.HasPrefix(s, "#") {
		return c, ExpectedHash
	}
	nibble := func(b byte) byte {
		switch {
		case b >= '0' && b <= '9':
			return b - '0'
		case b >= 'a' && b <= 'f':
			return b - 'a' + 10
		case b >= 'A' && b <= 'F':
			return b - 'A' + 10
		}
		err = InvalidPattern
		return 0
	}
	switch len(s) {
	case 9:
		c.R = float64(nibble(s[1])<<4+nibble(s[2])) / 255
		c.G = float64(nibble(s[3])<<4+nibble(s[4])) / 255
		c.B = float64(nibble(s[5])<<4+nibble(s[6])) / 255
		c.A = float64(nibble(s[7])<<4+nibble(s[8])) / 255
	case 7:
		c.R = float64(nibble(s[1])<<4+nibble(s[2])) / 255
		c.G = float64(nibble(s[3])<<4+nibble(s[4])) / 255
		c.B = float64(nibble(s[5])<<4+nibble(s[6])) / 255
		c.A = 1.0
	case 5:
		c.R = float64(nibble(s[1])*17) / 255
		c.G = float64(nibble(s[2])*17) / 255
		c.B = float64(nibble(s[3])*17) / 255
		c.A = float64(nibble(s[4])*17) / 255
	case 4:
		c.R = float64(nibble(s[1])*17) / 255
		c.G = float64(nibble(s[2])*17) / 255
		c.B = float64(nibble(s[3])*17) / 255
		c.A = 1.0
	default:
		err = InvalidLength
	}
	return c, err
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

// NRGBA returns the color as an NRGBA value.
func (c Color) NRGBA() color.NRGBA {
	return color.NRGBA{uint8(c.R * 255), uint8(c.G * 255), uint8(c.B * 255), uint8(c.A * 255)}
}

// Equal returns true when the colors have both equal color components as well as equal alpha value.
func (l Color) Equal(r Color) bool {
	return float.EqualPairs(l.R, r.R, l.G, r.G, l.B, r.B, l.A, r.A)
}

// Scale returns a Color with the rgb channels scaled by the supplied scalar value.
func (c Color) Scale(s float64) Color {
	c.R *= s
	c.G *= s
	c.B *= s
	return c
}

// AddChannels adds the channels of the current Color with each respective
// channel from the supplied Color object.
func (l Color) AddChannels(r Color) Color {
	l.R += r.R
	l.G += r.G
	l.B += r.B
	return l
}

// MultiplyChannels multiplies the channels of the current Color with each respective
// channel from the supplied Color object.
func (l Color) MultiplyChannels(r Color) Color {
	l.R *= r.R
	l.G *= r.G
	l.B *= r.B
	return l
}

// MinChannels takes the minimum between each channel of the current Color and the supplied Color object.
func (l Color) MinChannels(r Color) Color {
	l.R = math.Min(r.R, l.R)
	l.G = math.Min(r.G, l.G)
	l.B = math.Min(r.B, l.B)
	return l
}

// Clamp clamps each rgb channel to the supplied minimum and maximum scalar values.
func (l Color) Clamp(min, max float64) Color {
	l.R = math.Min(max, math.Max(min, l.R))
	l.G = math.Min(max, math.Max(min, l.G))
	l.B = math.Min(max, math.Max(min, l.B))
	return l
}
