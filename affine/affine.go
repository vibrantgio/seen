package affine

import (
	"github.com/reactivego/seen"
)

// This is the set of points that must be used by a surface that will use an
// affine transform for rendering.
// The coordinate system used by this basis is a right handed system where
// the x-axis is pointing right the y-axis is pointing up and the z axis
// is pointing out of the screen.
var ORTHONORMAL_BASIS = []seen.Point{
	{0, 0, 0},
	{20, 0, 0},
	{0, 20, 0},
}

// Matrix represents a transform in 2D used in SVG and HTML5 Canvas.
// The matrix can express rotation, skewing, scaling and translation.
// | A C E |
// | B D F |
// | 0 0 1 |
// An affine transformation respresented with a matrix uses homogeneous 
// coordinates. A column vector [x,y] is extended with a homogeneous 
// component w set to 1 like so: [x,y,1].
// Transforming the vector with the matrix will perform a 2x2 matrix
// operation with the top left 2x2 matrix [A,B,C,D] followed by adding 
// the translation [E,F] to produce [x',y',1]
type Matrix struct {
	A,B,C,D,E,F float64
}

// SolveForAffineTransform
// Computes the parameters of an affine transform from the 3 projected
// points. Points are specified in a coordinate system with positive x
// going right, positive y going up and positive z coming  out of the
// screen. Note that the Z component is not used for solving the affine
// transform.
// Returns affine transform matrix(A,B,C,D,E,F) interpreted as follows: 
//	| A C E |
//	| B D F |
//	| 0 0 1 |
// A scale, skew and rotation is the 2x2 matrix at the upper left and 
// a translation vector is at the upper right.
// NOTE! the coordinate system for homogeneous vectors that can be 
// transformed using this matrix have the positive x-axis going right
// but the positive y axis going down!
func SolveForAffineTransform(points []seen.Point) *Matrix {
	// Because we control the initial values of the points, we can re-use the
	// state matrix. Furthermore, because we have use a special layout (upper
	// triangular) for this matrix, we avoid any matrix factorization and can go
	// directly to back-substitution to solve the matrix equation.
	A := _INITIAL_STATE_MATRIX
	b := [...]float64{
		points[1].X,
		points[2].X,
		points[0].X,
		points[1].Y,
		points[2].Y,
		points[0].Y,
	}
	// Use back substitution to solve A*x=b for x
	var x [6]float64
	n := len(A)
	for i:=n-1; i>=0; i-- {
		x[i] = b[i]
		for j:=i+1; j<n; j++ {
			x[i] -= A[i][j] * x[j]
		}
		x[i] /= A[i][i]
	}
	// To use the affine transform, we flip y:
	//   x[0], x[3], -x[1], -x[4], x[2], x[5]
	return &Matrix{x[0],x[3],-x[1],-x[4],x[2],x[5]}
}

// _INITIAL_STATE_MATRIX is built using the method from this StackOverflow answer:
// http://stackoverflow.com/questions/22954239/given-three-points-compute-affine-transformation
// We further re-arranged the rows to avoid having to do any matrix factorization.
// The matrix consists of the ORTHONORMAL_BASIS vectors minus the Z component written in a
// different form that is appropriate for solving the affine transform.
var _INITIAL_STATE_MATRIX = [][]float64 {
	{20,  0, 1,  0,  0, 0},
	{ 0, 20, 1,  0,  0, 0},
	{ 0,  0, 1,  0,  0, 0},
	{ 0,  0, 0, 20,  0, 1},
	{ 0,  0, 0,  0, 20, 1},
	{ 0,  0, 0,  0,  0, 1},
}
