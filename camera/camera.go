package camera

import (
	"github.com/vibrantgio/seen/matrix"
	"github.com/vibrantgio/seen/projection"
	"github.com/vibrantgio/seen/transform"
)

// Camera contains all three major components of the 3D to 2D tranformation.
//
// First, we transform object from world-space (the same space that the coordinates of
// face points are in after all their transforms are applied) to camera space. Typically,
// this will place all viewable objects into a cube with coordinates:
// x = -1 to 1, y = -1 to 1, z = 1 to 2
//
// Second, we apply the projection transform to create perspective parallax and what not.
//
// Finally, we rescale to the viewport size.
//
// These three steps allow us to easily create shapes whose coordinates match up to
// screen coordinates in the z = 0 plane.
type Camera struct {
	transform.Transform
	Projection matrix.Matrix
}

var Default = CameraWithProjection(projection.DefaultPerspective)

func CameraWithProjection(projection matrix.Matrix) Camera {
	return Camera{Transform: transform.Default, Projection: projection}
}
