package seen

import "github.com/reactivego/seen/transform"

// Object base class extended by Shape and Model.
// Uses a double quaternion for specifying the transform.
type Object struct {
	dq transform.DualQuaternion
	sx float64
	sy float64
	sz float64
}

var DefaultObject = Object{transform.IdentDualQuaternion, 1.0, 1.0, 1.0}

// Matrix returns a 4x4 homogenous transformation matrix
// for the transform. This method makes Object a Transformable.
func (t *Object) Matrix() Matrix {
	m := t.dq.Mat4x4()
	if t.sx != 1.0 || t.sy != 1.0 || t.sz != 1.0 {
		m[0], m[1], m[2] = m[0]*t.sx, m[1]*t.sy, m[2]*t.sz
		m[4], m[5], m[6] = m[4]*t.sx, m[5]*t.sy, m[6]*t.sz
		m[8], m[9], m[10] = m[8]*t.sx, m[9]*t.sy, m[10]*t.sz
	}
	return Matrix{m}
}

// Rotation returns the Quaternion that specifies the rotation part of the transform.
func (t *Object) Rotation() transform.Quaternion {
	return t.dq.Rotation()
}

// SetRotation replaces the rotation part of the dual quaternion with a new rotation.
func (t *Object) SetRotation(r transform.Quaternion) {
	tx, ty, tz := t.dq.Translation()
	t.dq = transform.DualQuatRXYZ(r, tx, ty, tz)
}

// Translation returns the tx,ty,tz values that indicate the offset of the
// Object w.r.t. its parent object.
func (t *Object) Translation() (tx, ty, tz float64) {
	return t.dq.Translation()
}

// SetTranslation replaces the translation part of the dual quaternion with a new translation.
func (t *Object) SetTranslation(tx, ty, tz float64) {
	t.dq = transform.DualQuatRXYZ(t.dq.Rotation(), tx, ty, tz)
}

func (t *Object) Scale() (sx, sy, sz float64) {
	return t.sx, t.sy, t.sz
}

// SetScale sets the scaling to apply
func (t *Object) SetScale(sx, sy, sz float64) {
	t.sx, t.sy, t.sz = sx, sy, sz
}
