package viewport

import "github.com/vibrantgio/seen/matrix"

// Viewport
type Viewport struct{ Prescale, Postscale matrix.Matrix }

// Center creates a viewport in which the scene's origin maps to the centre of the
// view.
//
// By default the projection's scale and camera distance follow the view's width
// and height, so the scene fills the view and rescales as the view is resized.
//
// Pass an optional dist to lock the projection to that fixed reference distance
// instead. The scale and camera distance then no longer depend on width and
// height — those only position the centre — so one world unit projects to a
// constant number of pixels at any view size. Use this to keep on-screen size
// from changing as the window resizes; only the first dist value is used.
// Center(ox, oy, s, s) and Center(ox, oy, w, h, s) coincide when w == h == s.
func Center(offsetX, offsetY, width, height float64, dist ...float64) Viewport {
	// projW, projH, projD are the reference the projection is built against: the
	// view's own size by default, or the fixed dist when given. The centring
	// Translate below always uses the real width/height, never the reference.
	projW, projH, projD := width, height, height
	if len(dist) > 0 {
		projW, projH, projD = dist[0], dist[0], dist[0]
	}
	return Viewport{
		Prescale:  matrix.Scale(1/projW, 1/projH, 1/projD).Translate(-offsetX, -offsetY, -projD),
		Postscale: matrix.Translate(offsetX+width/2, offsetY+height/2, projD).Scale(projW, -projH, projD),
	}
}

// Origin creates a viewport in which the scene's origin aligns with the view's
// origin ([0, 0]), which is usually the top left.
//
// As with Center, the projection by default follows the view's width and height.
// Pass an optional dist to lock the scale and camera distance to that fixed
// reference instead, so on-screen size stays constant as the view is resized
// (width and height then only place the origin); only the first dist value is used.
func Origin(offsetX, offsetY, width, height float64, dist ...float64) Viewport {
	projW, projH, projD := width, height, height
	if len(dist) > 0 {
		projW, projH, projD = dist[0], dist[0], dist[0]
	}
	return Viewport{
		Prescale:  matrix.Scale(1/projW, 1/projH, 1/projD).Translate(-offsetX, -offsetY, -projD),
		Postscale: matrix.Translate(offsetX, offsetY, projD).Scale(projW, -projH, projD),
	}
}

var Default = Origin(0, 0, 1, 1)
