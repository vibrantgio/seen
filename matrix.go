package seen

import (
	"github.com/reactivego/seen/float"
	"github.com/reactivego/seen/transform"
)

// Matrix is a wrapper around transform.Matrix that can transform Point values.
type Matrix struct{ transform.Matrix }

var IdentityMatrix = Matrix{transform.IdentityMatrix}

func Translate(tx, ty, tz float64) Matrix {
	return Matrix{transform.Matrix{
		1, 0, 0, tx,
		0, 1, 0, ty,
		0, 0, 1, tz,
		0, 0, 0, 1,
	}}
}

func Scale(sx, sy, sz float64) Matrix {
	return Matrix{transform.Matrix{
		sx, 0, 0, 0,
		0, sy, 0, 0,
		0, 0, sz, 0,
		0, 0, 0, 1,
	}}
}

func (l Matrix) Mul(r Matrix) Matrix {
	return Matrix{l.Matrix.Mul(r.Matrix)}
}

func (l Matrix) Equal(r Matrix) bool {
	for i, li := range l.Matrix {
		if !float.Equal(li, r.Matrix[i]) {
			return false
		}
	}
	return true
}

func (m Matrix) Scale(sx, sy, sz float64) Matrix {
	s := transform.Matrix{
		sx, 0, 0, 0,
		0, sy, 0, 0,
		0, 0, sz, 0,
		0, 0, 0, 1,
	}
	return Matrix{m.Matrix.Mul(s)}
}

func (m Matrix) Translate(tx, ty, tz float64) Matrix {
	s := transform.Matrix{
		1, 0, 0, tx,
		0, 1, 0, ty,
		0, 0, 1, tz,
		0, 0, 0, 1,
	}
	return Matrix{m.Matrix.Mul(s)}
}
