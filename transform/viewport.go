package transform

// Viewport describes the screen coordinates to which projected coordinates can be converted.
// The view port origin is by default at the left top of the screen with the positive x axis
// going right and the positive y axis going down. The z axis can be choosen based on the
// depth range being specified. Projected space Z coordinate -1 representing the near plane
// can be mapped to a new value by specifying N in the viewport. A sensible default value
// for this is 0. The projected Z coordinate 1 representing the far plane can be mapped
// to a new value by specifying the F in the viewport. A sensible default value for this is 1.
// All projected Z coordinates in the range [-1,1] are mapped linear to [N,F].
type Viewport struct {
	// viewport
	X, Y, W, H float64

	// Depth range
	N, F float64
}

// Convert converts a projected coordinate to screen space.
func (v Viewport) Convert(x, y, z float64) (xs, ys, zs float64) {
	xs = v.X + x*0.5*v.W + 0.5*v.W
	ys = v.Y - y*0.5*v.H + 0.5*v.H
	zs = z*0.5*(v.F-v.N) + 0.5*(v.F+v.N)
	return
}
