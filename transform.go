package seen

import (
	"github.com/reactivego/seen/dualquat"
	"github.com/reactivego/seen/quat"
)

// Transformable is the interface every 3D object supports.
type Transformable interface {
	// Matrix returns the homogenous 4x4 matrix defining this Transformable's
	// coordinate system w.r.t. to its parent object.
	Matrix() Matrix
}

// Transform type is embedded by Light, Shape, Model and Camera.
// It uses a double quaternion for specifying the transform.
// The component transformations are applied in the order TRS.
// So Object coordinates are transformed into parent space by
// first scaling, then rotating and finally translating.
type Transform struct {
	dq dualquat.DualQuaternion
	sx float64
	sy float64
	sz float64
}

var DefaultTransform = Transform{dualquat.Identity, 1.0, 1.0, 1.0}

// Matrix returns a 4x4 homogenous transformation matrix for the transform.
// This method makes Transform a Transformable.
func (t *Transform) Matrix() Matrix {
	m := t.dq.Matrix()
	if t.sx != 1.0 || t.sy != 1.0 || t.sz != 1.0 {
		m[0][0], m[0][1], m[0][2] = m[0][0]*t.sx, m[0][1]*t.sy, m[0][2]*t.sz
		m[1][0], m[1][1], m[1][2] = m[1][0]*t.sx, m[1][1]*t.sy, m[1][2]*t.sz
		m[2][0], m[2][1], m[2][2] = m[2][0]*t.sx, m[2][1]*t.sy, m[2][2]*t.sz
	}
	return Matrix(m)
}

// Rotation returns the Quaternion that specifies the rotation part of the transform.
func (t *Transform) Rotation() quat.Quaternion {
	return t.dq.Rotation()
}

// SetRotation replaces the rotation part of the dual quaternion with a new rotation.
func (t *Transform) SetRotation(r quat.Quaternion) {
	tx, ty, tz := t.dq.Translation()
	t.dq = dualquat.TransRot(tx, ty, tz, r)
}

// Translation returns the tx,ty,tz values that indicate the offset of the
// Object w.r.t. its parent object.
func (t *Transform) Translation() (tx, ty, tz float64) {
	return t.dq.Translation()
}

// SetTranslation replaces the translation part of the dual quaternion with a new translation.
func (t *Transform) SetTranslation(tx, ty, tz float64) {
	t.dq = dualquat.TransRot(tx, ty, tz, t.dq.Rotation())
}

func (t *Transform) Scale() (sx, sy, sz float64) {
	return t.sx, t.sy, t.sz
}

// SetScale sets the scaling to apply
func (t *Transform) SetScale(sx, sy, sz float64) {
	t.sx, t.sy, t.sz = sx, sy, sz
}
