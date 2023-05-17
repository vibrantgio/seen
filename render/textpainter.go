package render

import "github.com/reactivego/seen/affine"

// TextPainter
type TextPainter interface {
	// FillText
	// transform is an affine matrix approximating a 3D transform of the plane on which the text is to be painted.
	// text is the text to be painted.
	// Style supports the following keys: fill, font, text-anchor
	FillText(transform affine.Matrix, text string, style Style)
}
