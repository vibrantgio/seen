package canvas

import (
	"github.com/vibrantgio/seen/affine"
	"github.com/vibrantgio/seen/point"
)

// Canvas
type Canvas interface {
	Path() PathPainter
	Rect() RectPainter
	Circle() CirclePainter
	Text() TextPainter

	Reset()
	Cleanup()
}

// PathPainter
type PathPainter interface {
	// Set up the path to be painted.
	// Then use Fill and/or Stroke to actually paint it.
	Path(points []point.Point) PathPainter

	// Fill the path
	Fill(Style)

	// Stroke the outline of the path.
	// Key "stroke-width" is supported in style.
	Stroke(Style)
}

// RectPainter
type RectPainter interface {
	Rect(width, height float64) RectPainter
	CornerRadius(rx, ry float64) RectPainter

	// Fill the rect
	Fill(Style)
}

// CirclePainter
type CirclePainter interface {
	Circle(center point.Point, radius float64) CirclePainter

	// Fill the circle
	Fill(Style)
}

// TextPainter
type TextPainter interface {
	// FillText
	// transform is an affine matrix approximating a 3D transform
	// 	of the plane on which the text is to be painted.
	// text is the text to be painted.
	// Style supports the following keys: "fill", "font", "text-anchor"
	FillText(transform affine.Matrix, text string, style Style)
}

type Style = map[string]string
