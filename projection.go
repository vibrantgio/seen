package seen

import "github.com/reactivego/seen/transform"

// DefaultPerspectiveProjection creates a perspective projection matrix from
// the default frustrum.
var DefaultPerspectiveProjection = Matrix{transform.Frustum(1, 1, 1, 100)}

// DefaultOrthographicProjection creates an orthographic projection matrix from
// the default frustrum.
var DefaultOrthographicProjection = Matrix{transform.Ortho(1, 1, 1, 100)}
