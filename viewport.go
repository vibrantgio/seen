package seen

// Viewport
type Viewport struct{ Prescale, Postscale Matrix }

// CenterViewport creates a viewport where the scene's origin is centered in the view
func CenterViewport(offsetX, offsetY, width, height float64) Viewport {
	return Viewport{
		Prescale:  Scale(1/width, 1/height, 1/height).Translate(-offsetX, -offsetY, -height),
		Postscale: Translate(offsetX+width/2, offsetY+height/2, height).Scale(width, -height, height),
	}
}

// OriginViewport creates a view port where the scene's origin is aligned with
// the origin ([0, 0]) of the view origin.
func OriginViewport(offsetX, offsetY, width, height float64) Viewport {
	return Viewport{
		Prescale:  Scale(1/width, 1/height, 1/height).Translate(-offsetX, -offsetY, -height),
		Postscale: Translate(offsetX, offsetY, height).Scale(width, -height, height),
	}
}

var DefaultViewport = OriginViewport(0, 0, 1, 1)
