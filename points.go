package seen

type Points []Point

func (points Points) Mul(m Matrix, transformedPoints Points) (barycenter Point) {
	if len(transformedPoints) != len(points) {
		panic("internal error, slice lengths don't match")
	}

	// Create Barycenter point used in sorting surfaces in the painters algorithm
	barycenter = Point{}

	// Apply transform to points
	for i, p := range points {
		p.X, p.Y, p.Z, _ = m.Transform(p.X, p.Y, p.Z, 1.0)
		transformedPoints[i] = p
		barycenter = barycenter.Add(p)
	}

	// Compute barycenter, which is used in sorting surfaces in the painters algorithm
	return barycenter.Scale(1.0 / float64(len(points)))
}

func (points Points) Clip(m Matrix, Z float64, clippedPoints Points) bool {
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
	v0 := p1.Subtract(p0)
	v1 := p2.Subtract(p0)
	return v0.Cross(v1)
}
