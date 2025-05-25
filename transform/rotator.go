package transform

// Rotator allows rotations to be compounded in a Transformer.
// For example, if T is a transform:
//
//	T.RotX(1.2).RotY(1.1).RotZ(0.9)
//
// Rotations are applied left to right on the vector being transformed: First
// the original rotation stored in the transform then RotX, then RotY, and
// finally RotZ.
type Rotator interface {
	RotX(angle float64) Rotator
	RotY(angle float64) Rotator
	RotZ(angle float64) Rotator
}
