package seen

import (
	"github.com/vibrantgio/seen/light"
	"github.com/vibrantgio/seen/matrix"
)

// Light represents a light in the scene graph.
type Light interface {
	Transformer
	IsEnabled() bool
	ShaderData(matrix.Matrix) light.ShaderData
}

var _ Light = (*light.Light)(nil)
