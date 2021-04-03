package seen

import "github.com/reactivego/seen/transform"

var defaultProjection = transform.Projection{R: 1, T: 1, N: 1, F: 100}

var DefaultPerspectiveProjection = Matrix{defaultProjection.PerspectiveMat4x4()}

var DefaultOrthographicProjection = Matrix{defaultProjection.OrthographicMat4x4()}
