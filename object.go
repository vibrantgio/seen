package seen

import (
	"github.com/reactivego/seen/transform"
)

// Object base class extended by Shape and Model.
// Uses a double quaternion for specifying the transform.
type Object struct {
	dq *transform.DualQuaternion
	scale *transform.Mat4x4
}

func (t *Object) Init() {
	t.dq = transform.IdentDualQuaternion
}

// Matrix returns a 4x4 homogenous transformation matrix
// for the transform. This method makes Object a Transformable.
func (t *Object) Matrix() *Matrix {
	if t.scale != nil {
		return &Matrix{t.dq.Mat3x4().Mat4x4().Mul(t.scale)}
	} else {
		return &Matrix{t.dq.Mat3x4().Mat4x4()}
	}
}

// Rotation returns the Quaternion that specifies the rotation part of the transform.
func (t *Object) Rotation() *transform.Quaternion {
	return t.dq.Rotation()
}

// SetRotation replaces the rotation part of the dual quaternion with a new rotation.
func (t *Object) SetRotation(r *transform.Quaternion) {
	tx, ty, tz := t.dq.Translation()
	t.dq = transform.MakeDualQuatRXYZ(r, tx, ty, tz)
}

// Translation returns the tx,ty,tz values that indicate the offset of the
// Object w.r.t. its parent object.
func (t *Object) Translation() (tx,ty,tz float64) {
	return t.dq.Translation()
}

// SetTranslation replaces the translation part of the dual quaternion with a new translation.
func (t *Object) SetTranslation(tx,ty,tz float64) {
	t.dq = transform.MakeDualQuatRXYZ(t.dq.Rotation(), tx, ty, tz)
}

func (t *Object) Scale() (sx,sy,sz float64) {
	if t.scale != nil {
		return t.scale[0], t.scale[5], t.scale[10]
	} else {
		return 1.0, 1.0, 1.0
	}
}

// SetScale sets the scaling to apply
func (t *Object) SetScale(sx,sy,sz float64) {
	t.scale = &transform.Mat4x4{
		sx,0,0,0,
		0,sy,0,0,
		0,0,sz,0,
		0,0,0,1,
	}
}
