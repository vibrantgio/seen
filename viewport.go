package seen

// Viewport
type Viewport struct {
	Prescale  *Matrix
	Postscale *Matrix
}

// MakeViewport
func MakeViewport() *Viewport {
	return &Viewport{IdentityMatrix, IdentityMatrix}
}

// MakeCenterViewport creates a viewport where the scene's origin is centered in the view
func MakeCenterViewport(offsetX, offsetY, width, height float64) *Viewport {
	x, y := offsetX, offsetY
	v := MakeViewport()
	v.Prescale = v.Prescale.Scale(1.0/width, 1.0/height, 1.0/height).Translate(-x, -y, -height)
	v.Postscale = v.Postscale.Translate(x + width/2.0, y + height/2.0, height).Scale(width, -height, height)
	return v
}

// MakeOriginViewport creates a view port where the scene's origin is aligned with
// the origin ([0, 0]) of the view origin.
func MakeOriginViewport(offsetX, offsetY, width, height float64) *Viewport {
	x, y := offsetX, offsetY
	v := MakeViewport()
	v.Prescale = v.Prescale.Scale(1.0/width, 1.0/height, 1.0/height).Translate(-x,-y,-1.0)
	v.Postscale = v.Postscale.Translate(x,y,0).Scale(width, -height, height)
	return v
}
