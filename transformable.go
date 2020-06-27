package seen

// Transformable is the interface every 3D object supports.
type Transformable interface {
	// Matrix returns the homogenous 4x4 matrix defining this Transformable's
	// coordinate system w.r.t. to its parent object.
	Matrix() *Matrix
}
