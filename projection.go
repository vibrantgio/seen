package seen

import "github.com/reactivego/seen/transform"

var defaultProjection = &transform.Projection{R: 1, T: 1, N: 1, F: 100}

func MakeDefaultPerspectiveProjection() *Matrix {
	return MakeMatrix(defaultProjection.PerspectiveMat4x4())
}

func MakeDefaultOrthographicProjection() *Matrix {
	return MakeMatrix(defaultProjection.OrthographicMat4x4())
}
