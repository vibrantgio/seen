package seen

// DefaultPerspectiveProjection creates a perspective projection matrix from
// the default frustrum.
var DefaultPerspectiveProjection = Frustum(1, 1, 1, 100)

// DefaultOrthographicProjection creates an orthographic projection matrix from
// the default frustrum.
var DefaultOrthographicProjection = Ortho(1, 1, 1, 100)
