package render

import "github.com/reactivego/seen"

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
