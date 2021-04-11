package seen

// Shape contains a collection of surfaces. They may create a closed 3D
// shape, but not necessarily. For example, a cube is a closed shape, but a
// patch is not.
type Shape struct {
	Type string
	Transform
	Surfaces
}
