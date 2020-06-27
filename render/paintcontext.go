package render

import (
	"github.com/reactivego/seen"
	"github.com/reactivego/seen/affine"
)

// PaintContext
type PaintContext interface {
	Path() PathPainter
	Rect() RectPainter
	Circle() CirclePainter
	Text() TextPainter

	Reset()
	Cleanup()
}

// PathPainter
type PathPainter interface {
	// Set up the path to be painted. Then use Fill and/or Stroke to 
	// actually paint it.
	Path(points []seen.Point)

	// Fill the path
	Fill(Style)

	// Stroke the outline of the path.
	// Key "stroke-width" is supported in style.
	Stroke(Style)
}

// RectPainter
type RectPainter interface {
	Size(width, height float64)
	CornerRadius(rx, ry float64)

	// Fill the rect
	Fill(Style)
}

// CirclePainter
type CirclePainter interface {
	Fill(Style)
}

// TextPainter
type TextPainter interface {
	// FillText
	// transform is an affine matrix approximating a 3D transform of the plane on which the text is to be painted.
	// text is the text to be painted.
	// Style supports the following keys: fill, font, text-anchor

	FillText(transform *affine.Matrix, text string, style Style)
}

type Style = map[string]string
