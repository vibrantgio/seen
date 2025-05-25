package viewport

import "github.com/vibrantgio/seen/matrix"

// Viewport
type Viewport struct{ Prescale, Postscale matrix.Matrix }

// Center creates a viewport where the scene's origin is centered in the view
func Center(offsetX, offsetY, width, height float64) Viewport {
	return Viewport{
		Prescale:  matrix.Scale(1/width, 1/height, 1/height).Translate(-offsetX, -offsetY, -height),
		Postscale: matrix.Translate(offsetX+width/2, offsetY+height/2, height).Scale(width, -height, height),
	}
}

// Origin creates a view port where the scene's origin is aligned with the
// origin ([0, 0]) of the view origin, which usually is at the top left.
func Origin(offsetX, offsetY, width, height float64) Viewport {
	return Viewport{
		Prescale:  matrix.Scale(1/width, 1/height, 1/height).Translate(-offsetX, -offsetY, -height),
		Postscale: matrix.Translate(offsetX, offsetY, height).Scale(width, -height, height),
	}
}

var Default = Origin(0, 0, 1, 1)
