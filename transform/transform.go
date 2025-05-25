package transform

import (
	"github.com/vibrantgio/seen/matrix"
	"github.com/vibrantgio/seen/quaternion"
)

// Transform is embedded by Light, Shape, Group, and Camera. It provides
// a transformation in 3D space, using a quaternion for rotation while
// translation and scaling remain separate. This approach enables flexible
// combination of these components as a single transformation. By following
// the T*R*S order (translation * rotation * scaling) from right to left,
// the object is scaled first, then rotated, and finally translated.
type Transform struct {
	sx, sy, sz float64
	r          quaternion.Quat
	tx, ty, tz float64
}

var Default = Transform{1.0, 1.0, 1.0, quaternion.Identity, 0.0, 0.0, 0.0}

// Rotation returns the Quaternion that specifies the rotation part of the transform.
func (t Transform) Rotation() quaternion.Quat {
	return t.r
}

// SetRotation replaces the rotation part of the dual quaternion with a new rotation.
func (t *Transform) SetRotation(r quaternion.Quat) {
	t.r = r
}

// Translation returns the tx,ty,tz values that indicate the offset of the
// Object w.r.t. its parent object.
func (t Transform) Translation() (tx, ty, tz float64) {
	return t.tx, t.ty, t.tz
}

// SetTranslation replaces the translation part of the dual quaternion with a new translation.
func (t *Transform) SetTranslation(tx, ty, tz float64) {
	t.tx, t.ty, t.tz = tx, ty, tz
}

func (t Transform) Scale() (sx, sy, sz float64) {
	return t.sx, t.sy, t.sz
}

// SetScale sets the scaling to apply
func (t *Transform) SetScale(sx, sy, sz float64) {
	t.sx, t.sy, t.sz = sx, sy, sz
}

func (t *Transform) RotX(angle float64) Rotator {
	t.SetRotation(quaternion.RotX(angle).Mul(t.Rotation()))
	return t
}

func (t *Transform) RotY(angle float64) Rotator {
	t.SetRotation(quaternion.RotY(angle).Mul(t.Rotation()))
	return t
}

func (t *Transform) RotZ(angle float64) Rotator {
	t.SetRotation(quaternion.RotZ(angle).Mul(t.Rotation()))
	return t
}

// Matrix returns a 4x4 homogenous transformation matrix for the transform.
func (t Transform) Matrix() matrix.Matrix {
	m := t.r.Matrix()
	if t.sx != 1.0 || t.sy != 1.0 || t.sz != 1.0 {
		m[0][0], m[0][1], m[0][2] = m[0][0]*t.sx, m[0][1]*t.sy, m[0][2]*t.sz
		m[1][0], m[1][1], m[1][2] = m[1][0]*t.sx, m[1][1]*t.sy, m[1][2]*t.sz
		m[2][0], m[2][1], m[2][2] = m[2][0]*t.sx, m[2][1]*t.sy, m[2][2]*t.sz
	}
	m[0][3], m[1][3], m[2][3] = t.tx, t.ty, t.tz
	return m
}

func (t Transform) Transform(vx, vy, vz float64) (x, y, z float64) {
	if t.sx != 1.0 || t.sy != 1.0 || t.sz != 1.0 {
		vx, vy, vz = t.sx*vx, t.sy*vy, t.sz*vz
	}
	vx, vy, vz = t.r.Transform(vx, vy, vz)
	return t.tx + vx, t.ty + vy, t.tz + vz
}
