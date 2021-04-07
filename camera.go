package seen

// Camera model contains all three major components of the 3D to 2D tranformation.
//
// First, we transform object from world-space (the same space that the coordinates of
// surface points are in after all their transforms are applied) to camera space. Typically,
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
	Object
	Projection Matrix
}

var DefaultCamera = CameraWithProjection(DefaultPerspectiveProjection)

func CameraWithProjection(projection Matrix) Camera {
	return Camera{Object: DefaultObject, Projection: projection}
}
