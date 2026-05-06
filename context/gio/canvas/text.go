package canvas

import (
	"fmt"
	"image"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/image/math/fixed"

	"gioui.org/f32"
	"gioui.org/font"
	"gioui.org/font/opentype"
	"gioui.org/io/system"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"

	"eliasnaur.com/font/roboto/robotobold"
	"eliasnaur.com/font/roboto/robotoregular"

	"github.com/vibrantgio/seen/affine"
	"github.com/vibrantgio/seen/canvas"
	"github.com/vibrantgio/seen/color"
)

// Locale

var Locale = EN_US
var EN_US = system.Locale{Language: "en-US", Direction: system.LTR}
var NL = system.Locale{Language: "nl", Direction: system.LTR}
var ZH_CN = system.Locale{Language: "zh-CN", Direction: system.LTR}

// Fonts

var RegularNormal = font.Font{Typeface: "Roboto", Style: font.Regular, Weight: font.Normal}
var RegularBold = font.Font{Typeface: "Roboto", Style: font.Regular, Weight: font.Bold}

var fontfaces struct {
	once       sync.Once
	collection []font.FontFace
}

func FontFaces() []font.FontFace {
	register := func(f font.Font, ttf []byte) {
		face, err := opentype.Parse(ttf)
		if err != nil {
			panic(fmt.Sprintf("failed to parse font: %v", err))
		}
		fontfaces.collection = append(fontfaces.collection, font.FontFace{Font: f, Face: face})
	}
	fontfaces.once.Do(func() {
		register(RegularNormal, robotoregular.TTF)
		register(RegularBold, robotobold.TTF)
	})
	n := len(fontfaces.collection)
	return fontfaces.collection[0:n:n]
}

// Text
type Text struct{ *op.Ops }

var shaper = text.NewShaper(text.WithCollection(FontFaces()))

// FillText
// transform is an affine matrix approximating a 3D transform of the plane on which the text is to be painted.
// text is the text to be painted.
// Style supports the following keys: fill, font, text-anchor
func (p *Text) FillText(t affine.Matrix, txt string, style canvas.Style) {
	aff := f32.NewAffine2D(float32(t.A), float32(t.C), float32(t.E), float32(t.B), float32(t.D), float32(t.F))
	defer op.Affine(aff).Push(p.Ops).Pop()

	fnt := font.Font{Typeface: "Roboto"}
	if family, present := style["font-family"]; present {
		fnt.Typeface = font.Typeface(family)
	}

	if weight, present := style["font-weight"]; present {
		switch weight {
		case "normal":
			fnt.Weight = font.Normal
		case "bold":
			fnt.Weight = font.Bold
		}
	}

	size := 10
	if sz, present := style["font-size"]; present {
		sz = strings.TrimSuffix(sz, "px")
		if sz, err := strconv.Atoi(sz); err == nil {
			size = sz
		}
	}

	fill := color.Black
	if c, present := style["fill"]; present {
		if f, err := color.ColorWithString(c); err == nil {
			fill = f
		}
	}

	ax, ay := float32(0.5), float32(1.0)
	if a, present := style["text-anchor"]; present {
		switch a {
		case "start":
			ax = 0.0
		case "middle":
			ax = 0.5
		case "end":
			ax = 1.0
		}
	}

	maxWidth := 2000
	if sz, present := style["inline-size"]; present {
		sz = strings.TrimSuffix(sz, "px")
		if sz, err := strconv.Atoi(sz); err == nil {
			maxWidth = sz
		}
	}

	// Layout the txt string given font, size and max width.
	params := text.Parameters{
		Font:     fnt,
		PxPerEm:  fixed.I(size),
		MaxWidth: maxWidth,
		Locale:   Locale,
	}
	shaper.LayoutString(params, txt)

	// Determine the size of the layout rectangle dx,dy
	dx, dy := 0, 0
	lines := [][]text.Glyph(nil)
	line := []text.Glyph(nil)
	for glyph, ok := shaper.NextGlyph(); ok; glyph, ok = shaper.NextGlyph() {
		line = append(line, glyph)
		if glyph.Flags&text.FlagLineBreak != 0 {
			dy += glyph.Ascent.Ceil() + glyph.Descent.Ceil()
			lineWidth := glyph.X.Ceil() + glyph.Advance.Ceil()
			if dx < lineWidth {
				dx = lineWidth
			}
			lines = append(lines, line)
			line = nil
		}
	}

	// Actually paint the txt
	offset := image.Pt(int(-ax*float32(dx)), int(-ay*float32(dy)))
	for _, line := range lines {
		shape := clip.Outline{Path: shaper.Shape(line)}.Op()
		glyph := line[len(line)-1]
		offset.Y += glyph.Ascent.Ceil()
		tstack := op.Offset(offset).Push(p.Ops)
		paint.FillShape(p.Ops, fill.NRGBA(), shape)
		offset.Y += glyph.Descent.Ceil()
		tstack.Pop()
	}
}
