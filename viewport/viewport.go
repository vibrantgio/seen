// Package viewport provides the screen mapping of the render pipeline: the
// matrix that takes normalized device coordinates (the result of the
// camera's projection and perspective divide) to pixels.
package viewport

import "github.com/vibrantgio/seen/matrix"

// Viewport maps normalized device coordinates to pixels. It carries no
// camera state — the eye position and view normalization live on the
// Camera; Screen is applied after projection and the perspective divide.
//
// Scene.FitCenter and Scene.FitOrigin configure Screen together with the
// camera's Eye and Norm to reproduce the legacy fill-the-view behaviour.
type Viewport struct{ Screen matrix.Matrix }

// Default maps the unit region with the scene origin at the view's origin
// — the screen mapping of the legacy Origin(0, 0, 1, 1) viewport. It pairs
// with camera.Default, which supplies the matching eye at (0, 0, 1).
var Default = Viewport{Screen: matrix.Translate(0, 0, 1).Scale(1, -1, 1)}
