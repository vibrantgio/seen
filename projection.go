package seen

import "github.com/reactivego/seen/transform"

var defaultProjection = transform.Projection{R: 1, T: 1, N: 1, F: 100}

// DefaultPerspectiveProjection creates a perspective projection matrix from
// the default frustrum.
var DefaultPerspectiveProjection = Matrix{defaultProjection.PerspectiveMat4x4()}

// DefaultOrthographicProjection creates an orthographic projection matrix from
// the default frustrum.
var DefaultOrthographicProjection = Matrix{defaultProjection.OrthographicMat4x4()}
