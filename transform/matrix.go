package transform

type Mat3x4 [12]float64

var IdentityMat3x4 = Mat3x4{
	1, 0, 0, 0,
	0, 1, 0, 0,
	0, 0, 1, 0,
}

type Mat4x4 [16]float64

var IdentityMat4x4 = Mat4x4{
	1, 0, 0, 0,
	0, 1, 0, 0,
	0, 0, 1, 0,
	0, 0, 0, 1,
}

// Multiply as if these are 4x4 matrixes with a bottom row of [0,0,0,1]
func (l Mat3x4) Mul(r Mat3x4) Mat3x4 {
	return Mat3x4{
		l[0]*r[0] + l[1]*r[4] + l[2]*r[8], l[0]*r[1] + l[1]*r[5] + l[2]*r[9], l[0]*r[2] + l[1]*r[6] + l[2]*r[10], l[0]*r[3] + l[1]*r[7] + l[2]*r[11] + l[3],
		l[4]*r[0] + l[5]*r[4] + l[6]*r[8], l[4]*r[1] + l[5]*r[5] + l[6]*r[9], l[4]*r[2] + l[5]*r[6] + l[6]*r[10], l[4]*r[3] + l[5]*r[7] + l[6]*r[11] + l[7],
		l[8]*r[0] + l[9]*r[4] + l[10]*r[8], l[8]*r[1] + l[9]*r[5] + l[10]*r[9], l[8]*r[2] + l[9]*r[6] + l[10]*r[10], l[8]*r[3] + l[9]*r[7] + l[10]*r[11] + l[11],
	}
}

// Transform transforms a point give in the matrix's object space coordinates into its parent space.
// To transform a single point, 9 muls and 9 adds are needed.
func (m Mat3x4) Transform(x, y, z float64) (rx, ry, rz float64) {
	rx = m[0]*x + m[1]*y + m[2]*z + m[3]
	ry = m[4]*x + m[5]*y + m[6]*z + m[7]
	rz = m[8]*x + m[9]*y + m[10]*z + m[11]
	return
}

// Return a 4x4 matrix with a bottom row of [0,0,0,1]
func (m Mat3x4) Mat4x4() Mat4x4 {
	return Mat4x4{
		m[0], m[1], m[2], m[3],
		m[4], m[5], m[6], m[7],
		m[8], m[9], m[10], m[11],
		0, 0, 0, 1,
	}
}

// Multiply proper 4x4 matrixes.
func (l Mat4x4) Mul(r Mat4x4) Mat4x4 {
	return Mat4x4{
		l[0]*r[0] + l[1]*r[4] + l[2]*r[8] + l[3]*r[12], l[0]*r[1] + l[1]*r[5] + l[2]*r[9] + l[3]*r[13], l[0]*r[2] + l[1]*r[6] + l[2]*r[10] + l[3]*r[14], l[0]*r[3] + l[1]*r[7] + l[2]*r[11] + l[3]*r[15],
		l[4]*r[0] + l[5]*r[4] + l[6]*r[8] + l[7]*r[12], l[4]*r[1] + l[5]*r[5] + l[6]*r[9] + l[7]*r[13], l[4]*r[2] + l[5]*r[6] + l[6]*r[10] + l[7]*r[14], l[4]*r[3] + l[5]*r[7] + l[6]*r[11] + l[7]*r[15],
		l[8]*r[0] + l[9]*r[4] + l[10]*r[8] + l[11]*r[12], l[8]*r[1] + l[9]*r[5] + l[10]*r[9] + l[11]*r[13], l[8]*r[2] + l[9]*r[6] + l[10]*r[10] + l[11]*r[14], l[8]*r[3] + l[9]*r[7] + l[10]*r[11] + l[11]*r[15],
		l[12]*r[0] + l[13]*r[4] + l[14]*r[8] + l[15]*r[12], l[12]*r[1] + l[13]*r[5] + l[14]*r[9] + l[15]*r[13], l[12]*r[2] + l[13]*r[6] + l[14]*r[10] + l[15]*r[14], l[12]*r[3] + l[13]*r[7] + l[14]*r[11] + l[15]*r[15],
	}
}

func (m Mat4x4) Transform(x, y, z, w float64) (rx, ry, rz, rw float64) {
	rx = m[0]*x + m[1]*y + m[2]*z + m[3]*w
	ry = m[4]*x + m[5]*y + m[6]*z + m[7]*w
	rz = m[8]*x + m[9]*y + m[10]*z + m[11]*w
	rw = m[12]*x + m[13]*y + m[14]*z + m[15]*w
	return
}
