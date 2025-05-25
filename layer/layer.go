package layer

import "github.com/vibrantgio/seen/canvas"

// Layer
type Layer interface {
	RenderOn(canvas.Canvas)
}
