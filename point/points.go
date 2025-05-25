package point

import "github.com/vibrantgio/seen/matrix"

type Points []Point

func (points Points) MulB(m matrix.Matrix, outPoints Points) (barycenter Point) {
	if len(outPoints) != len(points) {
		panic("internal error, slice lengths don't match")
	}

	// Apply transform to points
	for i, p := range points {
		p.X, p.Y, p.Z = m.Transform3(p.X, p.Y, p.Z)
		outPoints[i] = p
		barycenter = barycenter.Plus(p)
	}

	// Compute barycenter, which is used in sorting faces in the painters algorithm
	return barycenter.Times(1.0 / float64(len(points)))
}

func (points Points) Mul(m matrix.Matrix) Points {
	pts := make(Points, len(points))
	for i, p := range points {
		pts[i].X, pts[i].Y, pts[i].Z = m.Transform3(p.X, p.Y, p.Z)
	}
	return pts
}

// Barycenter of the points is used in sorting faces in the painters algorithm.
func (points Points) Barycenter() Point {
	var barycenter Point
	for _, p := range points {
		barycenter.X += p.X
		barycenter.Y += p.Y
		barycenter.Z += p.Z
	}
	return barycenter.Times(1.0 / float64(len(points)))
}

func (points Points) Clip(m matrix.Matrix, Z float64, clippedPoints Points) bool {
	if len(clippedPoints) != len(points) {
		panic("internal error, slice lengths don't match")
	}
	for i, p := range points {
		x, y, z, w := m.Transform(p.X, p.Y, p.Z, 1.0)
		if z <= Z {
			return false
		}
		// Apply the clip so it scales the x and y coordinates in a perspective projection.
		// This is done by dividing the X,Y,Z components by W
		clippedPoints[i] = Point{X: x / w, Y: y / w, Z: z / w}
	}
	return true
}

// Normal returns the normal (not normalize) for a slice
// of points. Can be used for e.g. backface culling or shading.
func (points Points) Normal() Point {
	pointsLen := len(points)
	if pointsLen < 3 {
		return Point{Z: 1} // Default normal
	}
	p0 := points[0]
	p1 := points[1]
	p2 := points[pointsLen-1]
	v0 := p1.Minus(p0)
	v1 := p2.Minus(p0)
	return v0.Cross(v1)
}

func (points Points) Round() {
	for i := range points {
		points[i] = points[i].Round()
	}
}
