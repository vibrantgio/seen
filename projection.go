package seen

import "github.com/reactivego/seen/dualquat"

// DefaultPerspectiveProjection creates a perspective projection matrix from
// the default frustrum.
var DefaultPerspectiveProjection = Matrix{dualquat.Frustum(1, 1, 1, 100)}

// DefaultOrthographicProjection creates an orthographic projection matrix from
// the default frustrum.
var DefaultOrthographicProjection = Matrix{dualquat.Ortho(1, 1, 1, 100)}
