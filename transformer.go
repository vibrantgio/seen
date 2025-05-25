package seen

import (
	"github.com/vibrantgio/seen/matrix"
	"github.com/vibrantgio/seen/quaternion"
	"github.com/vibrantgio/seen/transform"
)

// Transformer defines the interface for 3D objects that can be transformed in space.
// It provides methods to manipulate an object's position (translation), orientation (rotation),
// and size (scale), as well as access to its transformation matrix.
type Transformer interface {
	// Translation returns the object's position offset relative to its parent's coordinate system.
	// Returns tx, ty, tz as the displacement along the x, y, and z axes respectively.
	Translation() (tx, ty, tz float64)

	// SetTranslation updates the object's position in 3D space.
	// Parameters tx, ty, tz specify the displacement along the x, y, and z axes respectively.
	SetTranslation(tx, ty, tz float64)

	// Rotation returns the object's orientation as a quaternion.
	// The quaternion represents the rotation relative to the parent's coordinate system.
	Rotation() quaternion.Quat

	// SetRotation updates the object's orientation using the provided quaternion.
	// Parameter r specifies the new rotation to be applied.
	SetRotation(r quaternion.Quat)

	// Rotator introduces 3 methods that allow rotations to be specified more succintly.
	transform.Rotator

	// Scale returns the object's scaling factors along each axis.
	// Returns sx, sy, sz as the scale factors for x, y, and z axes respectively.
	Scale() (sx, sy, sz float64)

	// SetScale updates the object's size by applying scale factors along each axis.
	// Parameters sx, sy, sz specify the scale factors for x, y, and z axes respectively.
	SetScale(sx, sy, sz float64)

	// Matrix returns the complete transformation matrix for this object.
	// The 4x4 homogeneous matrix combines translation, rotation, and scale,
	// expressing the object's full transformation relative to its parent.
	Matrix() matrix.Matrix
}

var _ Transformer = (*transform.Transform)(nil)
